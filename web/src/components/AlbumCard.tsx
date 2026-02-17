import { Link } from 'react-router-dom';
import { Play } from 'lucide-react';
import type { Album } from '../api/types';
import { getArtworkUrl } from '../api/client';

interface AlbumCardProps {
  album: Album;
}

export function AlbumCard({ album }: AlbumCardProps) {
  return (
    <Link
      to={`/album/${album.id}`}
      className="group block rounded-md bg-mms-card p-3 hover:bg-mms-hover transition-colors"
    >
      {/* Cover art */}
      <div className="relative aspect-square rounded overflow-hidden bg-mms-surface mb-3">
        {album.cover_path ? (
          <img
            src={getArtworkUrl(album.id)}
            alt={album.title}
            className="w-full h-full object-cover"
            loading="lazy"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-mms-text-tertiary">
            <svg viewBox="0 0 24 24" className="w-12 h-12" fill="currentColor">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 14.5c-2.49 0-4.5-2.01-4.5-4.5S9.51 7.5 12 7.5s4.5 2.01 4.5 4.5-2.01 4.5-4.5 4.5zm0-5.5c-.55 0-1 .45-1 1s.45 1 1 1 1-.45 1-1-.45-1-1-1z" />
            </svg>
          </div>
        )}

        {/* Play overlay on hover */}
        <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
          <div className="w-10 h-10 rounded-full bg-mms-accent flex items-center justify-center shadow-lg">
            <Play size={20} fill="black" className="ml-0.5 text-black" />
          </div>
        </div>
      </div>

      {/* Info */}
      <p className="text-sm font-medium truncate">{album.title}</p>
      <p className="text-xs text-mms-text-secondary truncate mt-0.5">
        {album.artist_name}
        {album.year ? ` \u00B7 ${album.year}` : ''}
      </p>
    </Link>
  );
}
