package stream

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Streamer handles audio file streaming with Range request support.
type Streamer struct {
	cacheDir   string
	ffmpegPath string
}

// NewStreamer creates a new audio streamer.
func NewStreamer(cacheDir, ffmpegPath string) *Streamer {
	return &Streamer{
		cacheDir:   cacheDir,
		ffmpegPath: ffmpegPath,
	}
}

// ServeTrack streams an audio file with proper headers and Range support.
func (s *Streamer) ServeTrack(w http.ResponseWriter, r *http.Request, filePath, format string) {
	// Verify file exists
	if _, err := os.Stat(filePath); err != nil {
		log.Warn().Str("path", filePath).Msg("track file not found")
		http.Error(w, "track file not found", http.StatusNotFound)
		return
	}

	// Set content type based on format
	contentType := mimeTypeForFormat(format)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")

	// Use http.ServeFile which handles Range requests automatically
	http.ServeFile(w, r, filePath)
}

// mimeTypeForFormat returns the MIME type for an audio format.
func mimeTypeForFormat(format string) string {
	switch strings.ToLower(format) {
	case "flac":
		return "audio/flac"
	case "mp3":
		return "audio/mpeg"
	case "m4a", "aac":
		return "audio/mp4"
	case "ogg":
		return "audio/ogg"
	case "opus":
		return "audio/opus"
	case "wav":
		return "audio/wav"
	default:
		return "application/octet-stream"
	}
}

// TranscodeCachePath returns the cache path for a transcoded file.
func (s *Streamer) TranscodeCachePath(trackID, format string, bitrate int) string {
	filename := trackID + "_" + format + "_" + strings.Replace(
		filepath.Base(format), " ", "", -1) + "." + format
	return filepath.Join(s.cacheDir, filename)
}
