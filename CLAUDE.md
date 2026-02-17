# Mark's Music Solutions — Project Guide

Self-hosted music streaming server + ethical streaming research. Go backend (chi, SQLite FTS5, zerolog), React frontend (Zustand, Howler.js, Tailwind). Docker-ready, fully functional.

---

## Quick Reference

```bash
# Backend
cd server && go build -o mms ./cmd/mms/
./mms --config ../config.yaml              # http://localhost:8080
./mms --config ../config.yaml --scan       # Scan library on startup

# Frontend (dev)
cd web && npm install && npm run dev       # http://localhost:5173 (proxies API)

# Frontend (production build)
cd web && npm run build                    # Output: web/dist/

# Docker
docker compose up -d --build               # http://localhost:8080

# Lint/typecheck
cd web && npx tsc --noEmit
cd web && npx eslint .
```

---

## Architecture

```
marks-music-solutions/
├── CLAUDE.md
├── README.md                          # Research + ethical rankings
├── config.yaml                        # Server config
├── docker-compose.yml
├── research/                          # 7 streaming service analyses
├── data/
│   ├── mms.db                         # SQLite database (WAL mode)
│   └── artwork/                       # Extracted album covers
├── server/                            # Go backend
│   ├── cmd/mms/main.go               # Entry point (118 lines)
│   ├── internal/
│   │   ├── api/
│   │   │   ├── router.go             # chi/v5 routes (96 lines)
│   │   │   └── handlers.go           # HTTP handlers (345 lines)
│   │   ├── config/config.go          # YAML config
│   │   ├── db/
│   │   │   ├── db.go                 # SQLite + FTS5 setup (151 lines)
│   │   │   ├── models.go             # Data structures (99 lines)
│   │   │   └── repository.go         # Full data layer (532 lines)
│   │   ├── scanner/scanner.go        # Metadata extraction
│   │   ├── search/search.go          # FTS5 query prep
│   │   └── stream/stream.go          # Audio streaming + Range (70 lines)
│   ├── go.mod
│   └── Dockerfile                     # Multi-stage Alpine + FFmpeg
└── web/                               # React frontend
    ├── src/
    │   ├── api/ (client.ts, hooks.ts, types.ts)
    │   ├── components/ (Layout, NowPlayingBar, Sidebar, AlbumCard, TrackList)
    │   ├── pages/ (Home, Search, Library, Album, Artist)
    │   ├── store/player.ts           # Zustand + Howler.js (209 lines)
    │   └── hooks/useKeyboardShortcuts.ts
    ├── package.json
    ├── vite.config.ts                # API proxy → :8080
    └── tailwind.config.ts            # MMS dark theme
```

### API Endpoints (/api/v1)

```
GET  /artists, /artists/{id}, /artists/{id}/albums
GET  /albums, /albums/{id}, /albums/{id}/tracks
GET  /albums/recent?limit=20, /albums/random?limit=20
GET  /tracks/{id}, /tracks/{id}/stream (Range support)
GET  /artwork/{id}
GET  /search?q=query&limit=30 (FTS5 full-text)
POST /library/scan (202 Accepted, background)
CRUD /playlists, POST /playlists/{id}/tracks
POST /tracks/{id}/play (play history)
GET  /stats
```

---

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

---

## Conventions

- Go: chi router + handler struct, zerolog, context-aware repository
- React: client.ts → hooks.ts → component pattern
- Zustand wraps Howler.js for global player state
- FTS5 queries append `*` for prefix matching
- UPSERT pattern for idempotent metadata updates
- SPA fallback: all unmapped routes → index.html
- No authentication (local network assumption)

---

## Hard Assertions

```
NEVER: Block on audio streaming (write timeout disabled for Range requests)
ALWAYS: Use html5: true in Howler.js (streaming, not load-into-memory)
ALWAYS: Use WAL mode for SQLite
ALWAYS: Use parameterized SQL queries (no string interpolation)
ALWAYS: Record play history on track end
```

---

## Skills Index

| Skill | Scope | Use When |
|-------|-------|----------|
| `go-backend` | Project | chi router, SQLite FTS5, zerolog, streaming patterns |

---

## Cross-Repo Pointers

```
WORKSPACE:  ../CLAUDE.md
```

---

## Known Issues

- No authentication (local network only)
- CORS allows all origins (development mode)
- No test files exist yet
- Transcoding not yet implemented (cache path configured)
- Artist page partially complete, Settings page is stub
