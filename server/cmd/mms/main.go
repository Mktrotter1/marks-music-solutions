package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/marks-music-solutions/mms/internal/api"
	"github.com/marks-music-solutions/mms/internal/config"
	"github.com/marks-music-solutions/mms/internal/db"
	"github.com/marks-music-solutions/mms/internal/scanner"
	"github.com/marks-music-solutions/mms/internal/stream"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Structured logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Parse flags
	configPath := flag.String("config", "", "path to config.yaml")
	scanOnStart := flag.Bool("scan", false, "scan music library on startup")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadOrDefault(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	log.Info().
		Str("addr", cfg.Addr()).
		Strs("music_dirs", cfg.Music.Directories).
		Str("db", cfg.Database.Path).
		Msg("MMS starting")

	// Ensure data directories exist
	os.MkdirAll(filepath.Dir(cfg.Database.Path), 0755)
	os.MkdirAll(cfg.Transcode.CacheDir, 0755)
	os.MkdirAll("data/artwork", 0755)

	// Connect to database
	database, err := db.Open(cfg.Database.Path)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database")
	}
	defer database.Close()

	// Run migrations
	if err := db.Migrate(database); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	// Create repository
	repo := db.NewRepository(database)

	// Create scanner
	sc := scanner.NewScanner(repo, cfg.Music.Directories, "data/artwork")

	// Create streamer
	st := stream.NewStreamer(cfg.Transcode.CacheDir, cfg.Transcode.FFmpegPath)

	// Create handlers and router
	handlers := api.NewHandlers(repo, sc, st)
	router := api.NewRouter(handlers)

	// Scan on startup if requested
	if *scanOnStart {
		go func() {
			log.Info().Msg("starting initial library scan")
			if err := sc.ScanAll(); err != nil {
				log.Error().Err(err).Msg("initial scan failed")
			}
		}()
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      0, // No timeout for streaming responses
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info().Str("addr", cfg.Addr()).Msg("server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("shutting down gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("shutdown error")
	}

	log.Info().Msg("server stopped")
}
