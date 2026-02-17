import {
  Play,
  Pause,
  SkipBack,
  SkipForward,
  Volume2,
  VolumeX,
} from 'lucide-react';
import { usePlayerStore } from '../store/player';
import { getArtworkUrl } from '../api/client';
import { formatDuration } from '../lib/utils';
import { Link } from 'react-router-dom';

export function NowPlayingBar() {
  const {
    currentTrack,
    isPlaying,
    isLoading,
    volume,
    isMuted,
    currentTime,
    duration,
    togglePlay,
    next,
    previous,
    seek,
    setVolume,
    toggleMute,
  } = usePlayerStore();

  if (!currentTrack) return null;

  const progress = duration > 0 ? (currentTime / duration) * 100 : 0;

  return (
    <div className="h-player flex-shrink-0 bg-mms-surface border-t border-mms-border flex items-center px-4 gap-4">
      {/* Track info - left */}
      <div className="flex items-center gap-3 w-[240px] min-w-0">
        {/* Cover art */}
        <div className="w-12 h-12 rounded bg-mms-card flex-shrink-0 overflow-hidden">
          {currentTrack.cover_path ? (
            <img
              src={getArtworkUrl(currentTrack.album_id)}
              alt=""
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-mms-text-tertiary">
              <Volume2 size={20} />
            </div>
          )}
        </div>

        {/* Title + Artist */}
        <div className="min-w-0">
          <p className="text-sm font-medium truncate">{currentTrack.title}</p>
          <Link
            to={`/artist/${currentTrack.artist_id}`}
            className="text-xs text-mms-text-secondary hover:text-mms-accent truncate block"
          >
            {currentTrack.artist_name}
          </Link>
        </div>

        {/* Quality badge */}
        {currentTrack.bit_depth && currentTrack.sample_rate && (
          <span
            className={
              currentTrack.bit_depth > 16 || (currentTrack.sample_rate && currentTrack.sample_rate > 44100)
                ? 'quality-badge-hires'
                : 'quality-badge-lossless'
            }
          >
            {currentTrack.bit_depth}/{(currentTrack.sample_rate / 1000).toFixed(1)}
          </span>
        )}
      </div>

      {/* Player controls - center */}
      <div className="flex-1 flex flex-col items-center gap-1 max-w-[600px] mx-auto">
        {/* Buttons */}
        <div className="flex items-center gap-4">
          <button
            onClick={previous}
            className="text-mms-text-secondary hover:text-white transition-colors"
          >
            <SkipBack size={18} fill="currentColor" />
          </button>

          <button
            onClick={togglePlay}
            disabled={isLoading}
            className="w-8 h-8 rounded-full bg-white text-black flex items-center justify-center hover:scale-105 transition-transform disabled:opacity-50"
          >
            {isPlaying ? (
              <Pause size={16} fill="currentColor" />
            ) : (
              <Play size={16} fill="currentColor" className="ml-0.5" />
            )}
          </button>

          <button
            onClick={next}
            className="text-mms-text-secondary hover:text-white transition-colors"
          >
            <SkipForward size={18} fill="currentColor" />
          </button>
        </div>

        {/* Progress bar */}
        <div className="w-full flex items-center gap-2 text-[11px] text-mms-text-secondary">
          <span className="w-10 text-right tabular-nums">
            {formatDuration(currentTime)}
          </span>

          <div
            className="flex-1 h-1 bg-mms-border rounded-full cursor-pointer group relative"
            onClick={(e) => {
              const rect = e.currentTarget.getBoundingClientRect();
              const pct = (e.clientX - rect.left) / rect.width;
              seek(pct * duration);
            }}
          >
            <div
              className="h-full bg-white group-hover:bg-mms-accent rounded-full transition-colors relative"
              style={{ width: `${progress}%` }}
            >
              <div className="absolute right-0 top-1/2 -translate-y-1/2 w-3 h-3 bg-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity" />
            </div>
          </div>

          <span className="w-10 tabular-nums">
            {formatDuration(duration)}
          </span>
        </div>
      </div>

      {/* Volume - right */}
      <div className="flex items-center gap-2 w-[160px] justify-end">
        <button
          onClick={toggleMute}
          className="text-mms-text-secondary hover:text-white transition-colors"
        >
          {isMuted || volume === 0 ? <VolumeX size={16} /> : <Volume2 size={16} />}
        </button>
        <input
          type="range"
          min={0}
          max={1}
          step={0.01}
          value={isMuted ? 0 : volume}
          onChange={(e) => setVolume(Number(e.target.value))}
          className="w-24 accent-mms-accent"
        />
      </div>
    </div>
  );
}
