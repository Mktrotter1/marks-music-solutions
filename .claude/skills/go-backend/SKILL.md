# Skill: go-backend

Use when working on the Go music server: chi router, SQLite FTS5, zerolog, audio streaming, or metadata scanning.

## Router (chi/v5)

Middleware stack: RequestID → Logger → Recoverer → Timeout(60s) → CORS → Route handlers

### Routes (/api/v1)

```go
r.Get("/artists", h.ListArtists)           // Paginated, sorted by sort_name
r.Get("/artists/{id}", h.GetArtist)
r.Get("/artists/{id}/albums", h.GetArtistAlbums)
r.Get("/albums", h.ListAlbums)              // Paginated
r.Get("/albums/{id}", h.GetAlbum)
r.Get("/albums/{id}/tracks", h.ListTracksByAlbum)  // ORDER BY disc, track
r.Get("/albums/recent", h.RecentAlbums)     // ORDER BY created_at DESC
r.Get("/albums/random", h.RandomAlbums)     // ORDER BY RANDOM()
r.Get("/tracks/{id}", h.GetTrack)
r.Get("/tracks/{id}/stream", h.ServeTrack)  // Range support
r.Get("/artwork/{id}", h.ServeArtwork)      // Static file serve
r.Get("/search", h.Search)                  // FTS5 with prefix matching
r.Post("/library/scan", h.ScanLibrary)      // Background goroutine, 202 Accepted
r.Post("/tracks/{id}/play", h.RecordPlay)   // Play history
r.Get("/stats", h.GetStats)
```

SPA fallback: all unmapped routes → `web/dist/index.html`

## SQLite Configuration

```go
connStr = "data/mms.db?_journal_mode=WAL&_foreign_keys=on&_busy_timeout=5000"
db.SetMaxOpenConns(1)  // Single writer (SQLite limitation)
db.SetMaxIdleConns(1)
```

### Tables
- `artists` (id, name, sort_name, image_path)
- `albums` (id, artist_id, title, sort_title, year, genre, cover_path, track_count, disc_count, duration_seconds)
- `tracks` (id, album_id, artist_id, title, track_number, disc_number, duration_seconds, file_path, file_size, format, sample_rate, bit_depth, channels, bitrate)
- `playlists` + `playlist_tracks` (junction with position)
- `play_history` (track_id, played_at, duration_listened)
- `search_index` (FTS5 virtual table)
- `system_config` (KV store)

### FTS5 Search
- Tokenizer: `unicode61 remove_diacritics 2`
- Query: `query + "*"` for prefix matching (type-ahead)
- Indexed: title, artist, album
- Ranking: built-in FTS5 rank

## Repository Pattern

```go
type Repository struct { db *sql.DB }

// ID generation: SHA1-based deterministic
func generateID(parts ...string) string  // Same input = same ID

// Upsert: INSERT ... ON CONFLICT DO UPDATE (idempotent)
func (r *Repository) UpsertArtist(ctx, artist) error
func (r *Repository) UpsertAlbum(ctx, album) error
func (r *Repository) UpsertTrack(ctx, track) error

// After track upsert, update album stats
func (r *Repository) UpdateAlbumStats(ctx, albumID) error
```

## Audio Streaming

```go
// stream.go — Range request support via http.ServeFile()
func (s *Streamer) ServeTrack(w, r, track) {
    w.Header().Set("Content-Type", mimeType(track.Format))
    http.ServeFile(w, r, track.FilePath)  // Handles Range automatically
}

// MIME mapping
FLAC → audio/flac, MP3 → audio/mpeg, M4A → audio/mp4
WAV → audio/wav, OGG → audio/ogg
```

Write timeout MUST be disabled (0) for streaming to work.

## Metadata Scanner

- Libraries: `dhowden/tag` (ID3, MP4, Vorbis), `mewkiz/flac` (deep FLAC parsing)
- Extracts: title, artist, album, year, genre, track/disc numbers, duration, sample_rate, bit_depth, channels, bitrate
- Album art: extracted from tags, saved to `data/artwork/{albumID}.jpg`
- Thread-safe scanning (mutex)
- Background scan via goroutine (HTTP 202 Accepted)

## Pagination

- Offset-based (not cursor)
- Limit clamped: `max(1, min(limit, 200))`, default 50
- Response: `{ items: [], total: count }`
