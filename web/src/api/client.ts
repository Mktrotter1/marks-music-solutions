import type {
  Artist,
  Album,
  Track,
  Playlist,
  AlbumDetail,
  SearchResults,
  PaginatedResponse,
  LibraryStats,
} from './types';

const BASE = '/api/v1';

async function fetchJSON<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, init);
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return res.json();
}

// --- Artists ---

export function getArtists(limit = 50, offset = 0) {
  return fetchJSON<PaginatedResponse<Artist>>(
    `${BASE}/artists?limit=${limit}&offset=${offset}`,
  );
}

export function getArtist(id: string) {
  return fetchJSON<Artist>(`${BASE}/artists/${id}`);
}

export function getArtistAlbums(id: string) {
  return fetchJSON<PaginatedResponse<Album>>(`${BASE}/artists/${id}/albums`);
}

// --- Albums ---

export function getAlbums(limit = 50, offset = 0) {
  return fetchJSON<PaginatedResponse<Album>>(
    `${BASE}/albums?limit=${limit}&offset=${offset}`,
  );
}

export function getAlbum(id: string) {
  return fetchJSON<AlbumDetail>(`${BASE}/albums/${id}`);
}

export function getAlbumTracks(id: string) {
  return fetchJSON<PaginatedResponse<Track>>(`${BASE}/albums/${id}/tracks`);
}

export function getRecentAlbums(limit = 20) {
  return fetchJSON<PaginatedResponse<Album>>(
    `${BASE}/albums/recent?limit=${limit}`,
  );
}

export function getRandomAlbums(limit = 20) {
  return fetchJSON<PaginatedResponse<Album>>(
    `${BASE}/albums/random?limit=${limit}`,
  );
}

// --- Tracks ---

export function getTrack(id: string) {
  return fetchJSON<Track>(`${BASE}/tracks/${id}`);
}

export function getStreamUrl(trackId: string) {
  return `${BASE}/tracks/${trackId}/stream`;
}

// --- Search ---

export function search(query: string, limit = 30) {
  return fetchJSON<SearchResults>(
    `${BASE}/search?q=${encodeURIComponent(query)}&limit=${limit}`,
  );
}

// --- Artwork ---

export function getArtworkUrl(albumId: string) {
  return `${BASE}/artwork/${albumId}`;
}

// --- Playlists ---

export function getPlaylists() {
  return fetchJSON<PaginatedResponse<Playlist>>(`${BASE}/playlists`);
}

export function getPlaylist(id: string) {
  return fetchJSON<Playlist>(`${BASE}/playlists/${id}`);
}

export function createPlaylist(name: string, description = '') {
  return fetchJSON<Playlist>(`${BASE}/playlists`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, description }),
  });
}

export function deletePlaylist(id: string) {
  return fetch(`${BASE}/playlists/${id}`, { method: 'DELETE' });
}

export function addTrackToPlaylist(playlistId: string, trackId: string) {
  return fetch(`${BASE}/playlists/${playlistId}/tracks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ track_id: trackId }),
  });
}

// --- Library ---

export function scanLibrary() {
  return fetchJSON<{ status: string; message: string }>(`${BASE}/library/scan`, {
    method: 'POST',
  });
}

// --- Stats ---

export function getStats() {
  return fetchJSON<LibraryStats>(`${BASE}/stats`);
}

// --- Play tracking ---

export function recordPlay(trackId: string, duration?: number) {
  return fetch(`${BASE}/tracks/${trackId}/play`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ duration }),
  });
}
