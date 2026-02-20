## Key Constants

```
BACKEND:
  Framework:  go-chi/v5
  Database:   SQLite WAL + FTS5 (unicode61, remove_diacritics)
  Logging:    zerolog (structured JSON)
  Port:       8080 (default)
  Formats:    FLAC, MP3, M4A, OGG, OPUS, WAV
  Pagination: Default 50, max 200

FRONTEND:
  React 18, Zustand 5.0, TanStack React Query 5.62
  Howler.js 2.2.4 (html5: true), Tailwind 3.4, Vite 6.0
  TypeScript 5.6 strict

THEME: Pure black (#000) + cyan accent (#00FFFF)

KEYBOARD: Space=play/pause, ←/→=seek, ↑/↓=volume, M=mute, /=search

SQLITE: WAL mode, FKs on, 5s busy timeout, single writer
  Tables: artists, albums, tracks, playlists, playlist_tracks, play_history, search_index (FTS5)
  IDs: SHA1-based deterministic
```
