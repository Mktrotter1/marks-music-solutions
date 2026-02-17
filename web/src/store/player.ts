import { create } from 'zustand';
import { Howl } from 'howler';
import type { Track } from '../api/types';
import { getStreamUrl, recordPlay } from '../api/client';

interface PlayerState {
  // Current track
  currentTrack: Track | null;
  queue: Track[];
  queueIndex: number;

  // Playback state
  isPlaying: boolean;
  volume: number;
  currentTime: number;
  duration: number;
  isMuted: boolean;
  isLoading: boolean;

  // Howl instance (internal)
  _howl: Howl | null;
  _rafId: number | null;

  // Actions
  playTrack: (track: Track, queue?: Track[], queueIndex?: number) => void;
  playAlbum: (tracks: Track[], startIndex?: number) => void;
  togglePlay: () => void;
  pause: () => void;
  resume: () => void;
  stop: () => void;
  next: () => void;
  previous: () => void;
  seek: (time: number) => void;
  setVolume: (vol: number) => void;
  toggleMute: () => void;
  addToQueue: (track: Track) => void;
  clearQueue: () => void;
}

export const usePlayerStore = create<PlayerState>((set, get) => ({
  currentTrack: null,
  queue: [],
  queueIndex: -1,
  isPlaying: false,
  volume: 0.8,
  currentTime: 0,
  duration: 0,
  isMuted: false,
  isLoading: false,
  _howl: null,
  _rafId: null,

  playTrack: (track, queue, queueIndex) => {
    const state = get();

    // Stop current playback
    if (state._howl) {
      state._howl.unload();
    }
    if (state._rafId) {
      cancelAnimationFrame(state._rafId);
    }

    const howl = new Howl({
      src: [getStreamUrl(track.id)],
      html5: true, // Required for streaming large files
      volume: state.isMuted ? 0 : state.volume,
      onplay: () => {
        set({ isPlaying: true, isLoading: false });
        updateProgress();
      },
      onpause: () => set({ isPlaying: false }),
      onstop: () => set({ isPlaying: false, currentTime: 0 }),
      onend: () => {
        // Record play
        recordPlay(track.id, track.duration_seconds);
        // Auto-advance to next track
        get().next();
      },
      onload: () => {
        set({ duration: howl.duration(), isLoading: false });
      },
      onloaderror: (_id, error) => {
        console.error('Load error:', error);
        set({ isLoading: false });
      },
    });

    const updateProgress = () => {
      const seek = howl.seek();
      if (typeof seek === 'number') {
        set({ currentTime: seek });
      }
      if (howl.playing()) {
        const id = requestAnimationFrame(updateProgress);
        set({ _rafId: id });
      }
    };

    set({
      currentTrack: track,
      _howl: howl,
      isLoading: true,
      currentTime: 0,
      ...(queue !== undefined ? { queue, queueIndex: queueIndex ?? 0 } : {}),
    });

    howl.play();
  },

  playAlbum: (tracks, startIndex = 0) => {
    if (tracks.length === 0) return;
    get().playTrack(tracks[startIndex], tracks, startIndex);
  },

  togglePlay: () => {
    const { _howl, isPlaying } = get();
    if (!_howl) return;
    if (isPlaying) {
      _howl.pause();
    } else {
      _howl.play();
    }
  },

  pause: () => {
    get()._howl?.pause();
  },

  resume: () => {
    get()._howl?.play();
  },

  stop: () => {
    const state = get();
    if (state._howl) {
      state._howl.unload();
    }
    if (state._rafId) {
      cancelAnimationFrame(state._rafId);
    }
    set({
      currentTrack: null,
      isPlaying: false,
      currentTime: 0,
      duration: 0,
      _howl: null,
      _rafId: null,
    });
  },

  next: () => {
    const { queue, queueIndex } = get();
    if (queue.length === 0) return;
    const nextIndex = queueIndex + 1;
    if (nextIndex < queue.length) {
      get().playTrack(queue[nextIndex], queue, nextIndex);
    } else {
      get().stop();
    }
  },

  previous: () => {
    const { queue, queueIndex, currentTime } = get();
    // If more than 3 seconds in, restart current track
    if (currentTime > 3) {
      get().seek(0);
      return;
    }
    if (queue.length === 0) return;
    const prevIndex = queueIndex - 1;
    if (prevIndex >= 0) {
      get().playTrack(queue[prevIndex], queue, prevIndex);
    }
  },

  seek: (time) => {
    const { _howl } = get();
    if (_howl) {
      _howl.seek(time);
      set({ currentTime: time });
    }
  },

  setVolume: (vol) => {
    const { _howl, isMuted } = get();
    set({ volume: vol });
    if (_howl && !isMuted) {
      _howl.volume(vol);
    }
  },

  toggleMute: () => {
    const { _howl, isMuted, volume } = get();
    set({ isMuted: !isMuted });
    if (_howl) {
      _howl.volume(isMuted ? volume : 0);
    }
  },

  addToQueue: (track) => {
    set((state) => ({ queue: [...state.queue, track] }));
  },

  clearQueue: () => {
    set({ queue: [], queueIndex: -1 });
  },
}));
