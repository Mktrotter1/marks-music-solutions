import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as api from './client';

// --- Artists ---

export function useArtists(limit = 50, offset = 0) {
  return useQuery({
    queryKey: ['artists', limit, offset],
    queryFn: () => api.getArtists(limit, offset),
  });
}

export function useArtist(id: string) {
  return useQuery({
    queryKey: ['artist', id],
    queryFn: () => api.getArtist(id),
    enabled: !!id,
  });
}

export function useArtistAlbums(id: string) {
  return useQuery({
    queryKey: ['artist-albums', id],
    queryFn: () => api.getArtistAlbums(id),
    enabled: !!id,
  });
}

// --- Albums ---

export function useAlbums(limit = 50, offset = 0) {
  return useQuery({
    queryKey: ['albums', limit, offset],
    queryFn: () => api.getAlbums(limit, offset),
  });
}

export function useAlbum(id: string) {
  return useQuery({
    queryKey: ['album', id],
    queryFn: () => api.getAlbum(id),
    enabled: !!id,
  });
}

export function useRecentAlbums(limit = 20) {
  return useQuery({
    queryKey: ['albums-recent', limit],
    queryFn: () => api.getRecentAlbums(limit),
  });
}

export function useRandomAlbums(limit = 20) {
  return useQuery({
    queryKey: ['albums-random', limit],
    queryFn: () => api.getRandomAlbums(limit),
  });
}

// --- Search ---

export function useSearch(query: string, limit = 30) {
  return useQuery({
    queryKey: ['search', query, limit],
    queryFn: () => api.search(query, limit),
    enabled: query.length >= 2,
  });
}

// --- Playlists ---

export function usePlaylists() {
  return useQuery({
    queryKey: ['playlists'],
    queryFn: api.getPlaylists,
  });
}

export function usePlaylist(id: string) {
  return useQuery({
    queryKey: ['playlist', id],
    queryFn: () => api.getPlaylist(id),
    enabled: !!id,
  });
}

export function useCreatePlaylist() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ name, description }: { name: string; description?: string }) =>
      api.createPlaylist(name, description),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['playlists'] }),
  });
}

// --- Stats ---

export function useStats() {
  return useQuery({
    queryKey: ['stats'],
    queryFn: api.getStats,
  });
}

// --- Library scan ---

export function useScanLibrary() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: api.scanLibrary,
    onSuccess: () => {
      // Invalidate all library queries after scan
      setTimeout(() => {
        queryClient.invalidateQueries({ queryKey: ['albums'] });
        queryClient.invalidateQueries({ queryKey: ['artists'] });
        queryClient.invalidateQueries({ queryKey: ['stats'] });
        queryClient.invalidateQueries({ queryKey: ['albums-recent'] });
      }, 2000);
    },
  });
}
