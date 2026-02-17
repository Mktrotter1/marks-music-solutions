package scanner

import (
	"context"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dhowden/tag"
	"github.com/google/uuid"
	"github.com/marks-music-solutions/mms/internal/db"
	"github.com/mewkiz/flac"
	"github.com/rs/zerolog/log"
)

// Scanner walks music directories and extracts metadata into the database.
type Scanner struct {
	repo       *db.Repository
	dirs       []string
	artworkDir string
	mu         sync.Mutex
	scanning   bool
}

// NewScanner creates a new library scanner.
func NewScanner(repo *db.Repository, dirs []string, artworkDir string) *Scanner {
	return &Scanner{
		repo:       repo,
		dirs:       dirs,
		artworkDir: artworkDir,
	}
}

// IsScanning returns whether a scan is currently in progress.
func (s *Scanner) IsScanning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.scanning
}

// ScanAll walks all configured directories and indexes music files.
func (s *Scanner) ScanAll() error {
	s.mu.Lock()
	if s.scanning {
		s.mu.Unlock()
		return fmt.Errorf("scan already in progress")
	}
	s.scanning = true
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		s.scanning = false
		s.mu.Unlock()
	}()

	var total, errors int
	for _, dir := range s.dirs {
		n, e := s.scanDirectory(dir)
		total += n
		errors += e
	}

	log.Info().Int("tracks", total).Int("errors", errors).Msg("library scan complete")
	return nil
}

func (s *Scanner) scanDirectory(dir string) (scanned, errors int) {
	log.Info().Str("dir", dir).Msg("scanning directory")

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Warn().Err(err).Str("path", path).Msg("walk error")
			return nil // continue walking
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".flac" && ext != ".mp3" && ext != ".m4a" && ext != ".ogg" && ext != ".opus" {
			return nil
		}

		if err := s.scanFile(path, ext); err != nil {
			log.Warn().Err(err).Str("path", path).Msg("scan file error")
			errors++
		} else {
			scanned++
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Str("dir", dir).Msg("walk directory error")
	}
	return
}

func (s *Scanner) scanFile(path, ext string) error {
	ctx := context.Background()

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	// Extract metadata using dhowden/tag
	metadata, err := tag.ReadFrom(f)
	if err != nil {
		return fmt.Errorf("read metadata: %w", err)
	}

	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	// Get artist name (fall back to "Unknown Artist")
	artistName := metadata.Artist()
	if artistName == "" {
		artistName = metadata.AlbumArtist()
	}
	if artistName == "" {
		artistName = "Unknown Artist"
	}

	// Get album title (fall back to "Unknown Album")
	albumTitle := metadata.Album()
	if albumTitle == "" {
		albumTitle = "Unknown Album"
	}

	// Get track title (fall back to filename)
	title := metadata.Title()
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	// Upsert artist
	artist, err := s.repo.UpsertArtist(ctx, artistName, sortName(artistName))
	if err != nil {
		return fmt.Errorf("upsert artist: %w", err)
	}

	// Get year and genre
	var year *int
	if y := metadata.Year(); y != 0 {
		year = &y
	}
	var genre *string
	if g := metadata.Genre(); g != "" {
		genre = &g
	}

	// Upsert album
	album, err := s.repo.UpsertAlbum(ctx, artist.ID, albumTitle, sortName(albumTitle), year, genre)
	if err != nil {
		return fmt.Errorf("upsert album: %w", err)
	}

	// Get track/disc numbers
	trackNum, _ := metadata.Track()
	discNum, _ := metadata.Disc()
	if discNum == 0 {
		discNum = 1
	}

	// Get duration (FLAC-specific via mewkiz/flac for accuracy)
	var duration float64
	var sampleRate, bitDepth, channels *int
	if ext == ".flac" {
		duration, sampleRate, bitDepth, channels = s.getFLACInfo(path)
	}
	// Fallback: estimate from file size if duration is still 0
	if duration == 0 {
		// Rough estimate: assume ~1000 kbps for FLAC
		duration = float64(fi.Size()) / 125000
	}

	// Generate deterministic track ID from file path
	trackID := uuid.NewSHA1(uuid.NameSpaceURL, []byte("track:"+path)).String()

	format := strings.TrimPrefix(ext, ".")

	// Calculate bitrate
	var bitrate *int
	if duration > 0 {
		br := int(float64(fi.Size()*8) / duration)
		bitrate = &br
	}

	var trackNumPtr *int
	if trackNum > 0 {
		trackNumPtr = &trackNum
	}

	track := &db.Track{
		ID:              trackID,
		AlbumID:         album.ID,
		ArtistID:        artist.ID,
		Title:           title,
		TrackNumber:     trackNumPtr,
		DiscNumber:      discNum,
		DurationSeconds: duration,
		FilePath:        path,
		FileSize:        fi.Size(),
		Format:          format,
		SampleRate:      sampleRate,
		BitDepth:        bitDepth,
		Channels:        2,
		Bitrate:         bitrate,
	}
	if channels != nil {
		track.Channels = *channels
	}

	if err := s.repo.UpsertTrack(ctx, track); err != nil {
		return fmt.Errorf("upsert track: %w", err)
	}

	// Update album stats
	s.repo.UpdateAlbumStats(ctx, album.ID)

	// Extract cover art if album doesn't have one yet
	if album.CoverPath == nil {
		s.extractCoverArt(ctx, f, metadata, album.ID)
	}

	// Index for full-text search
	s.repo.IndexTrack(ctx, trackID, "track", title, artistName, albumTitle)
	s.repo.IndexTrack(ctx, album.ID, "album", albumTitle, artistName, "")
	s.repo.IndexTrack(ctx, artist.ID, "artist", artistName, "", "")

	return nil
}

// getFLACInfo reads FLAC STREAMINFO for accurate duration and audio properties.
func (s *Scanner) getFLACInfo(path string) (duration float64, sampleRate, bitDepth, channels *int) {
	stream, err := flac.ParseFile(path)
	if err != nil {
		return 0, nil, nil, nil
	}
	defer stream.Close()

	info := stream.Info
	if info.SampleRate > 0 && info.NSamples > 0 {
		duration = float64(info.NSamples) / float64(info.SampleRate)
	}

	sr := int(info.SampleRate)
	bd := int(info.BitsPerSample)
	ch := int(info.NChannels)
	return duration, &sr, &bd, &ch
}

// extractCoverArt saves embedded album art to disk.
func (s *Scanner) extractCoverArt(ctx context.Context, f *os.File, metadata tag.Metadata, albumID string) {
	pic := metadata.Picture()
	if pic == nil {
		return
	}

	// Ensure artwork directory exists
	os.MkdirAll(s.artworkDir, 0755)

	coverPath := filepath.Join(s.artworkDir, albumID+".jpg")

	// If it's already JPEG data, write directly
	if pic.MIMEType == "image/jpeg" {
		if err := os.WriteFile(coverPath, pic.Data, 0644); err != nil {
			log.Warn().Err(err).Str("album", albumID).Msg("failed to write cover art")
			return
		}
	} else {
		// For other formats, try to decode and re-encode as JPEG
		// For simplicity, just write the raw data with appropriate extension
		if pic.MIMEType == "image/png" {
			coverPath = filepath.Join(s.artworkDir, albumID+".png")
		}
		if err := os.WriteFile(coverPath, pic.Data, 0644); err != nil {
			log.Warn().Err(err).Str("album", albumID).Msg("failed to write cover art")
			return
		}
	}

	s.repo.UpdateAlbumCover(ctx, albumID, coverPath)
}

// sortName generates a sort-friendly name (strips leading "The ", etc.)
func sortName(name string) string {
	lower := strings.ToLower(name)
	for _, prefix := range []string{"the ", "a ", "an "} {
		if strings.HasPrefix(lower, prefix) {
			return name[len(prefix):]
		}
	}
	return name
}

// Ensure jpeg import is used (for potential future encoding)
var _ = jpeg.DefaultQuality
