export interface Artist {
  id: string;
  name: string;
  sort_name: string;
  image_path?: string;
  album_count?: number;
  track_count?: number;
  created_at: string;
  updated_at: string;
}

export interface Album {
  id: string;
  artist_id: string;
  title: string;
  sort_title: string;
  year?: number;
  genre?: string;
  cover_path?: string;
  track_count: number;
  disc_count: number;
  duration_seconds: number;
  artist_name?: string;
  created_at: string;
  updated_at: string;
}

export interface Track {
  id: string;
  album_id: string;
  artist_id: string;
  title: string;
  track_number?: number;
  disc_number: number;
  duration_seconds: number;
  file_size: number;
  format: string;
  sample_rate?: number;
  bit_depth?: number;
  channels: number;
  bitrate?: number;
  artist_name?: string;
  album_title?: string;
  cover_path?: string;
  created_at: string;
  updated_at: string;
}

export interface Playlist {
  id: string;
  name: string;
  description?: string;
  cover_path?: string;
  track_count: number;
  duration_seconds: number;
  created_at: string;
  updated_at: string;
}

export interface SearchResult {
  entity_id: string;
  entity_type: 'artist' | 'album' | 'track';
  title: string;
  artist: string;
  album: string;
  rank: number;
}

export interface SearchResults {
  artists: SearchResult[];
  albums: SearchResult[];
  tracks: SearchResult[];
  total: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
}

export interface AlbumDetail {
  album: Album;
  tracks: Track[];
}

export interface LibraryStats {
  artists: number;
  albums: number;
  tracks: number;
}
