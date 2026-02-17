package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Repository provides data access for the music library.
// All queries use parameterized statements.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new data repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// --- Artist Operations ---

// UpsertArtist creates or updates an artist by name.
func (r *Repository) UpsertArtist(ctx context.Context, name, sortName string) (*Artist, error) {
	id := uuid.NewSHA1(uuid.NameSpaceURL, []byte("artist:"+name)).String()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO artists (id, name, sort_name)
		 VALUES (?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   name = excluded.name,
		   updated_at = CURRENT_TIMESTAMP`,
		id, name, sortName,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert artist: %w", err)
	}
	return r.GetArtistByID(ctx, id)
}

// GetArtistByID retrieves an artist by ID.
func (r *Repository) GetArtistByID(ctx context.Context, id string) (*Artist, error) {
	a := &Artist{}
	err := r.db.QueryRowContext(ctx,
		`SELECT a.id, a.name, a.sort_name, a.image_path, a.created_at, a.updated_at,
		        (SELECT COUNT(*) FROM albums WHERE artist_id = a.id) as album_count,
		        (SELECT COUNT(*) FROM tracks WHERE artist_id = a.id) as track_count
		 FROM artists a WHERE a.id = ?`, id,
	).Scan(&a.ID, &a.Name, &a.SortName, &a.ImagePath, &a.CreatedAt, &a.UpdatedAt,
		&a.AlbumCount, &a.TrackCount)
	if err != nil {
		return nil, fmt.Errorf("get artist %s: %w", id, err)
	}
	return a, nil
}

// ListArtists returns all artists sorted by name.
func (r *Repository) ListArtists(ctx context.Context, limit, offset int) ([]*Artist, int64, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM artists`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count artists: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT a.id, a.name, a.sort_name, a.image_path, a.created_at, a.updated_at,
		        (SELECT COUNT(*) FROM albums WHERE artist_id = a.id) as album_count,
		        (SELECT COUNT(*) FROM tracks WHERE artist_id = a.id) as track_count
		 FROM artists a
		 ORDER BY a.sort_name ASC
		 LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list artists: %w", err)
	}
	defer rows.Close()

	var artists []*Artist
	for rows.Next() {
		a := &Artist{}
		if err := rows.Scan(&a.ID, &a.Name, &a.SortName, &a.ImagePath, &a.CreatedAt, &a.UpdatedAt,
			&a.AlbumCount, &a.TrackCount); err != nil {
			return nil, 0, fmt.Errorf("scan artist: %w", err)
		}
		artists = append(artists, a)
	}
	if artists == nil {
		artists = []*Artist{}
	}
	return artists, total, nil
}

// --- Album Operations ---

// UpsertAlbum creates or updates an album.
func (r *Repository) UpsertAlbum(ctx context.Context, artistID, title, sortTitle string, year *int, genre *string) (*Album, error) {
	id := uuid.NewSHA1(uuid.NameSpaceURL, []byte("album:"+artistID+":"+title)).String()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO albums (id, artist_id, title, sort_title, year, genre)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   title = excluded.title,
		   year = COALESCE(excluded.year, albums.year),
		   genre = COALESCE(excluded.genre, albums.genre),
		   updated_at = CURRENT_TIMESTAMP`,
		id, artistID, title, sortTitle, year, genre,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert album: %w", err)
	}
	return r.GetAlbumByID(ctx, id)
}

// GetAlbumByID retrieves an album by ID with artist name.
func (r *Repository) GetAlbumByID(ctx context.Context, id string) (*Album, error) {
	a := &Album{}
	err := r.db.QueryRowContext(ctx,
		`SELECT al.id, al.artist_id, al.title, al.sort_title, al.year, al.genre,
		        al.cover_path, al.track_count, al.disc_count, al.duration_seconds,
		        al.created_at, al.updated_at,
		        ar.name as artist_name
		 FROM albums al
		 JOIN artists ar ON ar.id = al.artist_id
		 WHERE al.id = ?`, id,
	).Scan(&a.ID, &a.ArtistID, &a.Title, &a.SortTitle, &a.Year, &a.Genre,
		&a.CoverPath, &a.TrackCount, &a.DiscCount, &a.DurationSeconds,
		&a.CreatedAt, &a.UpdatedAt, &a.ArtistName)
	if err != nil {
		return nil, fmt.Errorf("get album %s: %w", id, err)
	}
	return a, nil
}

// ListAlbums returns albums with pagination.
func (r *Repository) ListAlbums(ctx context.Context, limit, offset int) ([]*Album, int64, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM albums`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count albums: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT al.id, al.artist_id, al.title, al.sort_title, al.year, al.genre,
		        al.cover_path, al.track_count, al.disc_count, al.duration_seconds,
		        al.created_at, al.updated_at,
		        ar.name as artist_name
		 FROM albums al
		 JOIN artists ar ON ar.id = al.artist_id
		 ORDER BY al.sort_title ASC
		 LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list albums: %w", err)
	}
	defer rows.Close()

	return r.scanAlbums(rows)
}

// ListAlbumsByArtist returns all albums for a given artist.
func (r *Repository) ListAlbumsByArtist(ctx context.Context, artistID string) ([]*Album, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT al.id, al.artist_id, al.title, al.sort_title, al.year, al.genre,
		        al.cover_path, al.track_count, al.disc_count, al.duration_seconds,
		        al.created_at, al.updated_at,
		        ar.name as artist_name
		 FROM albums al
		 JOIN artists ar ON ar.id = al.artist_id
		 WHERE al.artist_id = ?
		 ORDER BY al.year DESC, al.sort_title ASC`, artistID,
	)
	if err != nil {
		return nil, fmt.Errorf("list albums by artist: %w", err)
	}
	defer rows.Close()

	albums, _, err := r.scanAlbums(rows)
	return albums, err
}

// RecentAlbums returns the most recently added albums.
func (r *Repository) RecentAlbums(ctx context.Context, limit int) ([]*Album, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT al.id, al.artist_id, al.title, al.sort_title, al.year, al.genre,
		        al.cover_path, al.track_count, al.disc_count, al.duration_seconds,
		        al.created_at, al.updated_at,
		        ar.name as artist_name
		 FROM albums al
		 JOIN artists ar ON ar.id = al.artist_id
		 ORDER BY al.created_at DESC
		 LIMIT ?`, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("recent albums: %w", err)
	}
	defer rows.Close()

	albums, _, err := r.scanAlbums(rows)
	return albums, err
}

// RandomAlbums returns random albums.
func (r *Repository) RandomAlbums(ctx context.Context, limit int) ([]*Album, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT al.id, al.artist_id, al.title, al.sort_title, al.year, al.genre,
		        al.cover_path, al.track_count, al.disc_count, al.duration_seconds,
		        al.created_at, al.updated_at,
		        ar.name as artist_name
		 FROM albums al
		 JOIN artists ar ON ar.id = al.artist_id
		 ORDER BY RANDOM()
		 LIMIT ?`, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("random albums: %w", err)
	}
	defer rows.Close()

	albums, _, err := r.scanAlbums(rows)
	return albums, err
}

func (r *Repository) scanAlbums(rows *sql.Rows) ([]*Album, int64, error) {
	var albums []*Album
	for rows.Next() {
		a := &Album{}
		if err := rows.Scan(&a.ID, &a.ArtistID, &a.Title, &a.SortTitle, &a.Year, &a.Genre,
			&a.CoverPath, &a.TrackCount, &a.DiscCount, &a.DurationSeconds,
			&a.CreatedAt, &a.UpdatedAt, &a.ArtistName); err != nil {
			return nil, 0, fmt.Errorf("scan album: %w", err)
		}
		albums = append(albums, a)
	}
	if albums == nil {
		albums = []*Album{}
	}
	return albums, int64(len(albums)), nil
}

// UpdateAlbumStats recalculates track_count, disc_count, and duration for an album.
func (r *Repository) UpdateAlbumStats(ctx context.Context, albumID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE albums SET
		   track_count = (SELECT COUNT(*) FROM tracks WHERE album_id = ?),
		   disc_count = (SELECT COALESCE(MAX(disc_number), 1) FROM tracks WHERE album_id = ?),
		   duration_seconds = (SELECT COALESCE(SUM(duration_seconds), 0) FROM tracks WHERE album_id = ?),
		   updated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`, albumID, albumID, albumID, albumID,
	)
	return err
}

// UpdateAlbumCover sets the cover art path for an album.
func (r *Repository) UpdateAlbumCover(ctx context.Context, albumID, coverPath string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE albums SET cover_path = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		coverPath, albumID,
	)
	return err
}

// --- Track Operations ---

// UpsertTrack creates or updates a track by file path.
func (r *Repository) UpsertTrack(ctx context.Context, t *Track) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tracks (id, album_id, artist_id, title, track_number, disc_number,
		                     duration_seconds, file_path, file_size, format,
		                     sample_rate, bit_depth, channels, bitrate)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   title = excluded.title,
		   track_number = excluded.track_number,
		   disc_number = excluded.disc_number,
		   duration_seconds = excluded.duration_seconds,
		   file_size = excluded.file_size,
		   sample_rate = excluded.sample_rate,
		   bit_depth = excluded.bit_depth,
		   channels = excluded.channels,
		   bitrate = excluded.bitrate,
		   updated_at = CURRENT_TIMESTAMP`,
		t.ID, t.AlbumID, t.ArtistID, t.Title, t.TrackNumber, t.DiscNumber,
		t.DurationSeconds, t.FilePath, t.FileSize, t.Format,
		t.SampleRate, t.BitDepth, t.Channels, t.Bitrate,
	)
	return err
}

// GetTrackByID retrieves a track with joined artist/album info.
func (r *Repository) GetTrackByID(ctx context.Context, id string) (*Track, error) {
	t := &Track{}
	err := r.db.QueryRowContext(ctx,
		`SELECT t.id, t.album_id, t.artist_id, t.title, t.track_number, t.disc_number,
		        t.duration_seconds, t.file_path, t.file_size, t.format,
		        t.sample_rate, t.bit_depth, t.channels, t.bitrate,
		        t.created_at, t.updated_at,
		        ar.name as artist_name, al.title as album_title, al.cover_path
		 FROM tracks t
		 JOIN artists ar ON ar.id = t.artist_id
		 JOIN albums al ON al.id = t.album_id
		 WHERE t.id = ?`, id,
	).Scan(&t.ID, &t.AlbumID, &t.ArtistID, &t.Title, &t.TrackNumber, &t.DiscNumber,
		&t.DurationSeconds, &t.FilePath, &t.FileSize, &t.Format,
		&t.SampleRate, &t.BitDepth, &t.Channels, &t.Bitrate,
		&t.CreatedAt, &t.UpdatedAt,
		&t.ArtistName, &t.AlbumTitle, &t.CoverPath)
	if err != nil {
		return nil, fmt.Errorf("get track %s: %w", id, err)
	}
	return t, nil
}

// ListTracksByAlbum returns all tracks for an album, ordered by disc/track number.
func (r *Repository) ListTracksByAlbum(ctx context.Context, albumID string) ([]*Track, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT t.id, t.album_id, t.artist_id, t.title, t.track_number, t.disc_number,
		        t.duration_seconds, t.file_path, t.file_size, t.format,
		        t.sample_rate, t.bit_depth, t.channels, t.bitrate,
		        t.created_at, t.updated_at,
		        ar.name as artist_name, al.title as album_title, al.cover_path
		 FROM tracks t
		 JOIN artists ar ON ar.id = t.artist_id
		 JOIN albums al ON al.id = t.album_id
		 WHERE t.album_id = ?
		 ORDER BY t.disc_number ASC, t.track_number ASC`, albumID,
	)
	if err != nil {
		return nil, fmt.Errorf("list tracks by album: %w", err)
	}
	defer rows.Close()

	return r.scanTracks(rows)
}

func (r *Repository) scanTracks(rows *sql.Rows) ([]*Track, error) {
	var tracks []*Track
	for rows.Next() {
		t := &Track{}
		if err := rows.Scan(&t.ID, &t.AlbumID, &t.ArtistID, &t.Title, &t.TrackNumber, &t.DiscNumber,
			&t.DurationSeconds, &t.FilePath, &t.FileSize, &t.Format,
			&t.SampleRate, &t.BitDepth, &t.Channels, &t.Bitrate,
			&t.CreatedAt, &t.UpdatedAt,
			&t.ArtistName, &t.AlbumTitle, &t.CoverPath); err != nil {
			return nil, fmt.Errorf("scan track: %w", err)
		}
		tracks = append(tracks, t)
	}
	if tracks == nil {
		tracks = []*Track{}
	}
	return tracks, nil
}

// --- Search Operations ---

// IndexTrack adds a track to the FTS5 search index.
func (r *Repository) IndexTrack(ctx context.Context, entityID, entityType, title, artist, album string) error {
	// Delete existing entry first (FTS5 doesn't support ON CONFLICT)
	r.db.ExecContext(ctx,
		`DELETE FROM search_index WHERE entity_id = ? AND entity_type = ?`,
		entityID, entityType)

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO search_index (entity_id, entity_type, title, artist, album)
		 VALUES (?, ?, ?, ?, ?)`,
		entityID, entityType, title, artist, album,
	)
	return err
}

// Search performs a full-text search across the library.
func (r *Repository) Search(ctx context.Context, query string, limit int) ([]*SearchResult, error) {
	if limit <= 0 || limit > 100 {
		limit = 30
	}

	// Use FTS5 match syntax with prefix matching
	ftsQuery := query + "*"

	rows, err := r.db.QueryContext(ctx,
		`SELECT entity_id, entity_type, title, artist, album, rank
		 FROM search_index
		 WHERE search_index MATCH ?
		 ORDER BY rank
		 LIMIT ?`, ftsQuery, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer rows.Close()

	var results []*SearchResult
	for rows.Next() {
		sr := &SearchResult{}
		if err := rows.Scan(&sr.EntityID, &sr.EntityType, &sr.Title, &sr.Artist, &sr.Album, &sr.Rank); err != nil {
			return nil, fmt.Errorf("scan search result: %w", err)
		}
		results = append(results, sr)
	}
	if results == nil {
		results = []*SearchResult{}
	}
	return results, nil
}

// --- Playlist Operations ---

// CreatePlaylist creates a new playlist.
func (r *Repository) CreatePlaylist(ctx context.Context, name, description string) (*Playlist, error) {
	id := uuid.New().String()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO playlists (id, name, description) VALUES (?, ?, ?)`,
		id, name, description,
	)
	if err != nil {
		return nil, fmt.Errorf("create playlist: %w", err)
	}
	return r.GetPlaylistByID(ctx, id)
}

// GetPlaylistByID retrieves a playlist.
func (r *Repository) GetPlaylistByID(ctx context.Context, id string) (*Playlist, error) {
	p := &Playlist{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, description, cover_path, track_count, duration_seconds,
		        created_at, updated_at
		 FROM playlists WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.CoverPath, &p.TrackCount,
		&p.DurationSeconds, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get playlist %s: %w", id, err)
	}
	return p, nil
}

// ListPlaylists returns all playlists.
func (r *Repository) ListPlaylists(ctx context.Context) ([]*Playlist, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, description, cover_path, track_count, duration_seconds,
		        created_at, updated_at
		 FROM playlists ORDER BY updated_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list playlists: %w", err)
	}
	defer rows.Close()

	var playlists []*Playlist
	for rows.Next() {
		p := &Playlist{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CoverPath, &p.TrackCount,
			&p.DurationSeconds, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan playlist: %w", err)
		}
		playlists = append(playlists, p)
	}
	if playlists == nil {
		playlists = []*Playlist{}
	}
	return playlists, nil
}

// DeletePlaylist removes a playlist.
func (r *Repository) DeletePlaylist(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM playlists WHERE id = ?`, id)
	return err
}

// AddTrackToPlaylist appends a track to a playlist.
func (r *Repository) AddTrackToPlaylist(ctx context.Context, playlistID, trackID string) error {
	id := uuid.New().String()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO playlist_tracks (id, playlist_id, track_id, position)
		 VALUES (?, ?, ?, (SELECT COALESCE(MAX(position), 0) + 1 FROM playlist_tracks WHERE playlist_id = ?))`,
		id, playlistID, trackID, playlistID,
	)
	if err != nil {
		return fmt.Errorf("add track to playlist: %w", err)
	}

	// Update playlist stats
	_, err = r.db.ExecContext(ctx,
		`UPDATE playlists SET
		   track_count = (SELECT COUNT(*) FROM playlist_tracks WHERE playlist_id = ?),
		   duration_seconds = (SELECT COALESCE(SUM(t.duration_seconds), 0) FROM playlist_tracks pt JOIN tracks t ON t.id = pt.track_id WHERE pt.playlist_id = ?),
		   updated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`, playlistID, playlistID, playlistID,
	)
	return err
}

// RemoveTrackFromPlaylist removes a track from a playlist by position.
func (r *Repository) RemoveTrackFromPlaylist(ctx context.Context, playlistID string, position int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM playlist_tracks WHERE playlist_id = ? AND position = ?`,
		playlistID, position,
	)
	return err
}

// --- Play History ---

// RecordPlay records a play event.
func (r *Repository) RecordPlay(ctx context.Context, trackID string, durationListened *float64) error {
	id := uuid.New().String()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO play_history (id, track_id, duration_listened) VALUES (?, ?, ?)`,
		id, trackID, durationListened,
	)
	return err
}

// --- Stats ---

// CountEntities returns total counts of artists, albums, and tracks.
func (r *Repository) CountEntities(ctx context.Context, artists, albums, tracks *int64) {
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM artists`).Scan(artists)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM albums`).Scan(albums)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tracks`).Scan(tracks)
}
