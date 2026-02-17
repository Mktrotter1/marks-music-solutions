package db

import "time"

// Artist represents a music artist.
type Artist struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	SortName  string    `json:"sort_name"`
	ImagePath *string   `json:"image_path,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Aggregated fields (not stored directly)
	AlbumCount int `json:"album_count,omitempty"`
	TrackCount int `json:"track_count,omitempty"`
}

// Album represents a music album.
type Album struct {
	ID              string    `json:"id"`
	ArtistID        string    `json:"artist_id"`
	Title           string    `json:"title"`
	SortTitle       string    `json:"sort_title"`
	Year            *int      `json:"year,omitempty"`
	Genre           *string   `json:"genre,omitempty"`
	CoverPath       *string   `json:"cover_path,omitempty"`
	TrackCount      int       `json:"track_count"`
	DiscCount       int       `json:"disc_count"`
	DurationSeconds float64   `json:"duration_seconds"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	// Joined fields
	ArtistName string `json:"artist_name,omitempty"`
}

// Track represents a music track.
type Track struct {
	ID              string    `json:"id"`
	AlbumID         string    `json:"album_id"`
	ArtistID        string    `json:"artist_id"`
	Title           string    `json:"title"`
	TrackNumber     *int      `json:"track_number,omitempty"`
	DiscNumber      int       `json:"disc_number"`
	DurationSeconds float64   `json:"duration_seconds"`
	FilePath        string    `json:"-"`
	FileSize        int64     `json:"file_size"`
	Format          string    `json:"format"`
	SampleRate      *int      `json:"sample_rate,omitempty"`
	BitDepth        *int      `json:"bit_depth,omitempty"`
	Channels        int       `json:"channels"`
	Bitrate         *int      `json:"bitrate,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	// Joined fields
	ArtistName string `json:"artist_name,omitempty"`
	AlbumTitle string `json:"album_title,omitempty"`
	CoverPath  *string `json:"cover_path,omitempty"`
}

// Playlist represents a user playlist.
type Playlist struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	CoverPath       *string   `json:"cover_path,omitempty"`
	TrackCount      int       `json:"track_count"`
	DurationSeconds float64   `json:"duration_seconds"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PlaylistTrack represents a track within a playlist.
type PlaylistTrack struct {
	ID         string    `json:"id"`
	PlaylistID string    `json:"playlist_id"`
	TrackID    string    `json:"track_id"`
	Position   int       `json:"position"`
	AddedAt    time.Time `json:"added_at"`
	// Joined track fields
	Track *Track `json:"track,omitempty"`
}

// PlayHistory represents a play event.
type PlayHistory struct {
	ID               string    `json:"id"`
	TrackID          string    `json:"track_id"`
	PlayedAt         time.Time `json:"played_at"`
	DurationListened *float64  `json:"duration_listened,omitempty"`
}

// SearchResult represents a full-text search result.
type SearchResult struct {
	EntityID   string  `json:"entity_id"`
	EntityType string  `json:"entity_type"` // "artist", "album", "track"
	Title      string  `json:"title"`
	Artist     string  `json:"artist"`
	Album      string  `json:"album"`
	Rank       float64 `json:"rank"`
}
