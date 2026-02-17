import { useParams, Link } from 'react-router-dom';
import { Play } from 'lucide-react';
import { useAlbum } from '../api/hooks';
import { getArtworkUrl } from '../api/client';
import { usePlayerStore } from '../store/player';
import { TrackList } from '../components/TrackList';
import { formatTotalDuration, formatQuality, pluralize } from '../lib/utils';

export function AlbumPage() {
  const { id } = useParams<{ id: string }>();
  const { data, isLoading } = useAlbum(id!);
  const playAlbum = usePlayerStore((s) => s.playAlbum);

  if (isLoading) {
    return <div className="text-mms-text-secondary">Loading...</div>;
  }

  if (!data) {
    return <div className="text-mms-text-secondary">Album not found</div>;
  }

  const { album, tracks } = data;

  // Determine highest quality track for quality badge
  const maxBitDepth = Math.max(...tracks.map((t) => t.bit_depth ?? 16));
  const maxSampleRate = Math.max(...tracks.map((t) => t.sample_rate ?? 44100));

  return (
    <div>
      {/* Album header */}
      <div className="flex gap-6 mb-8">
        {/* Cover art */}
        <div className="w-[200px] h-[200px] md:w-[240px] md:h-[240px] rounded-md overflow-hidden bg-mms-surface flex-shrink-0 shadow-2xl">
          {album.cover_path ? (
            <img
              src={getArtworkUrl(album.id)}
              alt={album.title}
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-mms-text-tertiary">
              <svg viewBox="0 0 24 24" className="w-20 h-20" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 14.5c-2.49 0-4.5-2.01-4.5-4.5S9.51 7.5 12 7.5s4.5 2.01 4.5 4.5-2.01 4.5-4.5 4.5zm0-5.5c-.55 0-1 .45-1 1s.45 1 1 1 1-.45 1-1-.45-1-1-1z" />
              </svg>
            </div>
          )}
        </div>

        {/* Album info */}
        <div className="flex flex-col justify-end min-w-0">
          <p className="text-xs uppercase tracking-wider text-mms-text-secondary mb-1">
            Album
          </p>
          <h1 className="text-3xl md:text-4xl font-bold mb-2 truncate">
            {album.title}
          </h1>
          <div className="flex items-center gap-2 text-sm text-mms-text-secondary flex-wrap">
            <Link
              to={`/artist/${album.artist_id}`}
              className="font-semibold text-white hover:text-mms-accent transition-colors"
            >
              {album.artist_name}
            </Link>
            {album.year && <span>&middot; {album.year}</span>}
            <span>
              &middot; {album.track_count} {pluralize(album.track_count, 'track')}
            </span>
            <span>&middot; {formatTotalDuration(album.duration_seconds)}</span>
          </div>

          {/* Quality badge */}
          <div className="mt-2">
            <span
              className={
                maxBitDepth > 16 || maxSampleRate > 44100
                  ? 'quality-badge-hires'
                  : 'quality-badge-lossless'
              }
            >
              {formatQuality(maxBitDepth, maxSampleRate, tracks[0]?.format)}
            </span>
          </div>

          {/* Play button */}
          <div className="mt-4">
            <button
              onClick={() => playAlbum(tracks)}
              className="inline-flex items-center gap-2 px-6 py-2 rounded-full bg-mms-accent text-black font-semibold text-sm hover:brightness-110 transition"
            >
              <Play size={18} fill="currentColor" />
              Play
            </button>
          </div>
        </div>
      </div>

      {/* Track list */}
      <TrackList tracks={tracks} showArtist={false} showNumber={true} />
    </div>
  );
}
