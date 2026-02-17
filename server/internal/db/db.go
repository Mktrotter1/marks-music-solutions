package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

// Open creates a new SQLite connection with WAL mode and FTS5 support.
func Open(path string) (*sql.DB, error) {
	// Enable WAL mode and foreign keys via connection string
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=on&_busy_timeout=5000", path)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	// Connection pool settings for SQLite (single-writer)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	log.Info().Str("path", path).Msg("database connected")
	return db, nil
}

// Migrate runs all schema migrations.
func Migrate(db *sql.DB) error {
	for i, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration %d failed: %w", i, err)
		}
	}
	log.Info().Int("count", len(migrations)).Msg("database migrations applied")
	return nil
}

var migrations = []string{
	// Artists table
	`CREATE TABLE IF NOT EXISTS artists (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		sort_name TEXT NOT NULL,
		image_path TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`,
	`CREATE INDEX IF NOT EXISTS idx_artists_name ON artists(name)`,
	`CREATE INDEX IF NOT EXISTS idx_artists_sort_name ON artists(sort_name)`,

	// Albums table
	`CREATE TABLE IF NOT EXISTS albums (
		id TEXT PRIMARY KEY,
		artist_id TEXT NOT NULL REFERENCES artists(id),
		title TEXT NOT NULL,
		sort_title TEXT NOT NULL,
		year INTEGER,
		genre TEXT,
		cover_path TEXT,
		track_count INTEGER DEFAULT 0,
		disc_count INTEGER DEFAULT 1,
		duration_seconds REAL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`,
	`CREATE INDEX IF NOT EXISTS idx_albums_artist_id ON albums(artist_id)`,
	`CREATE INDEX IF NOT EXISTS idx_albums_title ON albums(title)`,
	`CREATE INDEX IF NOT EXISTS idx_albums_year ON albums(year)`,
	`CREATE INDEX IF NOT EXISTS idx_albums_created_at ON albums(created_at)`,

	// Tracks table
	`CREATE TABLE IF NOT EXISTS tracks (
		id TEXT PRIMARY KEY,
		album_id TEXT NOT NULL REFERENCES albums(id),
		artist_id TEXT NOT NULL REFERENCES artists(id),
		title TEXT NOT NULL,
		track_number INTEGER,
		disc_number INTEGER DEFAULT 1,
		duration_seconds REAL NOT NULL,
		file_path TEXT NOT NULL UNIQUE,
		file_size INTEGER NOT NULL,
		format TEXT NOT NULL DEFAULT 'flac',
		sample_rate INTEGER,
		bit_depth INTEGER,
		channels INTEGER DEFAULT 2,
		bitrate INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`,
	`CREATE INDEX IF NOT EXISTS idx_tracks_album_id ON tracks(album_id)`,
	`CREATE INDEX IF NOT EXISTS idx_tracks_artist_id ON tracks(artist_id)`,
	`CREATE INDEX IF NOT EXISTS idx_tracks_file_path ON tracks(file_path)`,

	// Playlists table
	`CREATE TABLE IF NOT EXISTS playlists (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		cover_path TEXT,
		track_count INTEGER DEFAULT 0,
		duration_seconds REAL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`,

	// Playlist tracks junction table
	`CREATE TABLE IF NOT EXISTS playlist_tracks (
		id TEXT PRIMARY KEY,
		playlist_id TEXT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
		track_id TEXT NOT NULL REFERENCES tracks(id),
		position INTEGER NOT NULL,
		added_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`,
	`CREATE INDEX IF NOT EXISTS idx_playlist_tracks_playlist ON playlist_tracks(playlist_id, position)`,

	// Play history
	`CREATE TABLE IF NOT EXISTS play_history (
		id TEXT PRIMARY KEY,
		track_id TEXT NOT NULL REFERENCES tracks(id),
		played_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		duration_listened REAL
	)`,
	`CREATE INDEX IF NOT EXISTS idx_play_history_track ON play_history(track_id)`,
	`CREATE INDEX IF NOT EXISTS idx_play_history_played_at ON play_history(played_at)`,

	// FTS5 search index
	`CREATE VIRTUAL TABLE IF NOT EXISTS search_index USING fts5(
		entity_id UNINDEXED,
		entity_type UNINDEXED,
		title,
		artist,
		album,
		content='',
		tokenize='unicode61 remove_diacritics 2'
	)`,

	// System config
	`CREATE TABLE IF NOT EXISTS system_config (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`,
}
