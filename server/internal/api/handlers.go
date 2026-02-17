package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/marks-music-solutions/mms/internal/db"
	"github.com/marks-music-solutions/mms/internal/scanner"
	"github.com/marks-music-solutions/mms/internal/stream"
	"github.com/rs/zerolog/log"
)

// Handlers holds all HTTP handler dependencies.
type Handlers struct {
	repo     *db.Repository
	scanner  *scanner.Scanner
	streamer *stream.Streamer
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(repo *db.Repository, sc *scanner.Scanner, st *stream.Streamer) *Handlers {
	return &Handlers{
		repo:     repo,
		scanner:  sc,
		streamer: st,
	}
}

// --- Health ---

func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "mms",
	})
}

// --- Artists ---

func (h *Handlers) HandleListArtists(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	artists, total, err := h.repo.ListArtists(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list artists")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": artists,
		"total": total,
	})
}

func (h *Handlers) HandleGetArtist(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	artist, err := h.repo.GetArtistByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "artist not found")
		return
	}
	writeJSON(w, http.StatusOK, artist)
}

func (h *Handlers) HandleGetArtistAlbums(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	albums, err := h.repo.ListAlbumsByArtist(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list albums")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": albums,
		"total": len(albums),
	})
}

// --- Albums ---

func (h *Handlers) HandleListAlbums(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	albums, total, err := h.repo.ListAlbums(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list albums")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": albums,
		"total": total,
	})
}

func (h *Handlers) HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	album, err := h.repo.GetAlbumByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "album not found")
		return
	}

	// Also fetch tracks for the album
	tracks, err := h.repo.ListTracksByAlbum(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list tracks")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"album":  album,
		"tracks": tracks,
	})
}

func (h *Handlers) HandleGetAlbumTracks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tracks, err := h.repo.ListTracksByAlbum(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list tracks")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": tracks,
		"total": len(tracks),
	})
}

func (h *Handlers) HandleRecentAlbums(w http.ResponseWriter, r *http.Request) {
	limit := parseIntParam(r, "limit", 20)
	albums, err := h.repo.RecentAlbums(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get recent albums")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": albums,
		"total": len(albums),
	})
}

func (h *Handlers) HandleRandomAlbums(w http.ResponseWriter, r *http.Request) {
	limit := parseIntParam(r, "limit", 20)
	albums, err := h.repo.RandomAlbums(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get random albums")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": albums,
		"total": len(albums),
	})
}

// --- Tracks ---

func (h *Handlers) HandleGetTrack(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	track, err := h.repo.GetTrackByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "track not found")
		return
	}
	writeJSON(w, http.StatusOK, track)
}

func (h *Handlers) HandleStreamTrack(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	track, err := h.repo.GetTrackByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "track not found")
		return
	}
	h.streamer.ServeTrack(w, r, track.FilePath, track.Format)
}

// --- Artwork ---

func (h *Handlers) HandleArtwork(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	album, err := h.repo.GetAlbumByID(r.Context(), id)
	if err != nil || album.CoverPath == nil {
		writeError(w, http.StatusNotFound, "artwork not found")
		return
	}
	http.ServeFile(w, r, *album.CoverPath)
}

// --- Search ---

func (h *Handlers) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	limit := parseIntParam(r, "limit", 30)
	results, err := h.repo.Search(r.Context(), query, limit)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("search failed")
		writeError(w, http.StatusInternalServerError, "search failed")
		return
	}

	// Group results by type
	var artists, albums, tracks []*db.SearchResult
	for _, sr := range results {
		switch sr.EntityType {
		case "artist":
			artists = append(artists, sr)
		case "album":
			albums = append(albums, sr)
		case "track":
			tracks = append(tracks, sr)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"artists": artists,
		"albums":  albums,
		"tracks":  tracks,
		"total":   len(results),
	})
}

// --- Library Management ---

func (h *Handlers) HandleScanLibrary(w http.ResponseWriter, r *http.Request) {
	go func() {
		if err := h.scanner.ScanAll(); err != nil {
			log.Error().Err(err).Msg("library scan failed")
		}
	}()
	writeJSON(w, http.StatusAccepted, map[string]string{
		"status":  "scanning",
		"message": "Library scan started in background",
	})
}

// --- Playlists ---

func (h *Handlers) HandleListPlaylists(w http.ResponseWriter, r *http.Request) {
	playlists, err := h.repo.ListPlaylists(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list playlists")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": playlists,
		"total": len(playlists),
	})
}

func (h *Handlers) HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	playlist, err := h.repo.CreatePlaylist(r.Context(), body.Name, body.Description)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create playlist")
		return
	}
	writeJSON(w, http.StatusCreated, playlist)
}

func (h *Handlers) HandleGetPlaylist(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	playlist, err := h.repo.GetPlaylistByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "playlist not found")
		return
	}
	writeJSON(w, http.StatusOK, playlist)
}

func (h *Handlers) HandleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeletePlaylist(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete playlist")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) HandleAddTrackToPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "id")
	var body struct {
		TrackID string `json:"track_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.TrackID == "" {
		writeError(w, http.StatusBadRequest, "track_id is required")
		return
	}

	if err := h.repo.AddTrackToPlaylist(r.Context(), playlistID, body.TrackID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add track")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Play History ---

func (h *Handlers) HandleRecordPlay(w http.ResponseWriter, r *http.Request) {
	trackID := chi.URLParam(r, "id")
	var body struct {
		Duration *float64 `json:"duration"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if err := h.repo.RecordPlay(r.Context(), trackID, body.Duration); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to record play")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Stats ---

func (h *Handlers) HandleStats(w http.ResponseWriter, r *http.Request) {
	// Quick library stats
	ctx := r.Context()
	var artistCount, albumCount, trackCount int64
	h.repo.CountEntities(ctx, &artistCount, &albumCount, &trackCount)

	writeJSON(w, http.StatusOK, map[string]any{
		"artists": artistCount,
		"albums":  albumCount,
		"tracks":  trackCount,
	})
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func parsePagination(r *http.Request) (limit, offset int) {
	limit = parseIntParam(r, "limit", 50)
	offset = parseIntParam(r, "offset", 0)
	return
}

func parseIntParam(r *http.Request, name string, defaultVal int) int {
	s := r.URL.Query().Get(name)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}
