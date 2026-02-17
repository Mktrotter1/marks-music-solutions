// Format seconds to MM:SS or H:MM:SS
export function formatDuration(seconds: number): string {
  if (!seconds || isNaN(seconds)) return '0:00';

  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = Math.floor(seconds % 60);

  if (h > 0) {
    return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  }
  return `${m}:${s.toString().padStart(2, '0')}`;
}

// Format total duration for albums/playlists (e.g., "1 hr 23 min")
export function formatTotalDuration(seconds: number): string {
  if (!seconds) return '';
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  if (h > 0) {
    return `${h} hr ${m} min`;
  }
  return `${m} min`;
}

// Format file size
export function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`;
}

// Format quality string (e.g., "24bit/96kHz FLAC")
export function formatQuality(
  bitDepth?: number,
  sampleRate?: number,
  format?: string,
): string {
  const parts: string[] = [];
  if (bitDepth && sampleRate) {
    parts.push(`${bitDepth}bit/${(sampleRate / 1000).toFixed(sampleRate % 1000 === 0 ? 0 : 1)}kHz`);
  }
  if (format) {
    parts.push(format.toUpperCase());
  }
  return parts.join(' ');
}

// Pluralize a word
export function pluralize(count: number, singular: string, plural?: string): string {
  return count === 1 ? singular : (plural ?? singular + 's');
}
