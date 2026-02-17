import { Play, Pause } from 'lucide-react';
import type { Track } from '../api/types';
import { usePlayerStore } from '../store/player';
import { formatDuration, formatQuality } from '../lib/utils';

interface TrackListProps {
  tracks: Track[];
  showAlbum?: boolean;
  showArtist?: boolean;
  showNumber?: boolean;
}

export function TrackList({
  tracks,
  showAlbum = false,
  showArtist = true,
  showNumber = true,
}: TrackListProps) {
  const { currentTrack, isPlaying, playTrack } = usePlayerStore();

  const handlePlay = (track: Track, index: number) => {
    playTrack(track, tracks, index);
  };

  return (
    <div className="w-full">
      {/* Header row */}
      <div className="flex items-center gap-4 px-4 py-2 text-[11px] uppercase tracking-wider text-mms-text-tertiary border-b border-mms-border">
        <div className="w-8 text-center">#</div>
        <div className="flex-1">Title</div>
        {showArtist && <div className="w-[200px] hidden md:block">Artist</div>}
        {showAlbum && <div className="w-[200px] hidden lg:block">Album</div>}
        <div className="w-[80px] hidden sm:block text-right">Quality</div>
        <div className="w-[60px] text-right">Time</div>
      </div>

      {/* Track rows */}
      {tracks.map((track, index) => {
        const isCurrent = currentTrack?.id === track.id;
        return (
          <button
            key={track.id}
            onClick={() => handlePlay(track, index)}
            className={`w-full flex items-center gap-4 px-4 py-2 text-sm hover:bg-mms-hover transition-colors group ${
              isCurrent ? 'text-mms-accent' : ''
            }`}
          >
            {/* Track number / play icon */}
            <div className="w-8 text-center">
              <span className="group-hover:hidden">
                {isCurrent && isPlaying ? (
                  <Pause size={14} className="inline text-mms-accent" />
                ) : showNumber ? (
                  <span className="text-mms-text-secondary tabular-nums">
                    {track.track_number ?? index + 1}
                  </span>
                ) : (
                  <span className="text-mms-text-secondary tabular-nums">
                    {index + 1}
                  </span>
                )}
              </span>
              <Play
                size={14}
                className="hidden group-hover:inline text-white"
                fill="currentColor"
              />
            </div>

            {/* Title */}
            <div className="flex-1 text-left truncate">
              <span className={isCurrent ? 'text-mms-accent' : 'text-white'}>
                {track.title}
              </span>
            </div>

            {/* Artist */}
            {showArtist && (
              <div className="w-[200px] hidden md:block text-mms-text-secondary truncate text-left">
                {track.artist_name}
              </div>
            )}

            {/* Album */}
            {showAlbum && (
              <div className="w-[200px] hidden lg:block text-mms-text-secondary truncate text-left">
                {track.album_title}
              </div>
            )}

            {/* Quality */}
            <div className="w-[80px] hidden sm:block text-right">
              {track.bit_depth && track.sample_rate && (
                <span
                  className={
                    track.bit_depth > 16 || (track.sample_rate && track.sample_rate > 44100)
                      ? 'quality-badge-hires'
                      : 'quality-badge-lossless'
                  }
                >
                  {formatQuality(track.bit_depth, track.sample_rate)}
                </span>
              )}
            </div>

            {/* Duration */}
            <div className="w-[60px] text-right text-mms-text-secondary tabular-nums">
              {formatDuration(track.duration_seconds)}
            </div>
          </button>
        );
      })}
    </div>
  );
}
