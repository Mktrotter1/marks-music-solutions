package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates the HTTP router with middleware.
func NewRouter(handlers *Handlers) http.Handler {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(corsMiddleware)

	// Health checks
	r.Get("/healthz", handlers.HandleHealth)

	// Serve frontend static files
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/dist/assets"))))
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/dist/favicon.ico")
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Artists
		r.Get("/artists", handlers.HandleListArtists)
		r.Get("/artists/{id}", handlers.HandleGetArtist)
		r.Get("/artists/{id}/albums", handlers.HandleGetArtistAlbums)

		// Albums
		r.Get("/albums", handlers.HandleListAlbums)
		r.Get("/albums/{id}", handlers.HandleGetAlbum)
		r.Get("/albums/{id}/tracks", handlers.HandleGetAlbumTracks)
		r.Get("/albums/recent", handlers.HandleRecentAlbums)
		r.Get("/albums/random", handlers.HandleRandomAlbums)

		// Tracks
		r.Get("/tracks/{id}", handlers.HandleGetTrack)
		r.Get("/tracks/{id}/stream", handlers.HandleStreamTrack)

		// Artwork
		r.Get("/artwork/{id}", handlers.HandleArtwork)

		// Search
		r.Get("/search", handlers.HandleSearch)

		// Library management
		r.Post("/library/scan", handlers.HandleScanLibrary)

		// Playlists
		r.Get("/playlists", handlers.HandleListPlaylists)
		r.Post("/playlists", handlers.HandleCreatePlaylist)
		r.Get("/playlists/{id}", handlers.HandleGetPlaylist)
		r.Delete("/playlists/{id}", handlers.HandleDeletePlaylist)
		r.Post("/playlists/{id}/tracks", handlers.HandleAddTrackToPlaylist)

		// Play history
		r.Post("/tracks/{id}/play", handlers.HandleRecordPlay)

		// Stats
		r.Get("/stats", handlers.HandleStats)
	})

	// SPA fallback - serve index.html for all other routes
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/dist/index.html")
	})

	return r
}

// corsMiddleware adds CORS headers for development.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Range")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Range, Content-Length, Accept-Ranges")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
