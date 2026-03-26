# Mark's Music Solutions — Project Guide

Ethical streaming research + active Qobuz setup + self-hosted music server prototype.

## Active Music Stack

**Primary streaming: Qobuz via QBZ** (native Linux client, AUR `qbz-bin 1.1.18-1`)
- Qobuz user ID: 10620517
- QBZ is a Tauri (Rust/WebKit) app — all state in SQLite databases
- Full QBZ reference: `ref/qbz-client.md`

Key QBZ data paths:
```
~/.config/qbz/.qbz-auth              # Auth credentials (NEVER commit)
~/.local/share/qbz/                   # App data + databases
~/.local/share/qbz/users/10620517/    # User-specific data
~/.local/share/qbz/radio/             # Radio engine state
~/.cache/qbz/                         # Artwork + playback cache (~800MB)
```

QBZ integrations: Last.fm scrobbling, ListenBrainz, MusicBrainz, Discogs artwork, Chromecast, DLNA, AirPlay, Plex.

## Self-Hosted Server (MMS)

Go backend (chi, SQLite FTS5, zerolog), React frontend (Zustand, Howler.js, Tailwind). Docker-ready.

```bash
cd server && go build -o mms ./cmd/mms/    # Build backend
./mms --config ../config.yaml              # http://localhost:8080
./mms --config ../config.yaml --scan       # Scan library on startup
cd web && npm install && npm run dev       # http://localhost:5173 (proxies API)
cd web && npm run build                    # Output: web/dist/
docker compose up -d --build               # http://localhost:8080
cd web && npx tsc --noEmit && npx eslint . # Lint + typecheck
```

## Architecture
```
marks-music-solutions/
├── config.yaml / docker-compose.yml / README.md
├── research/                          # 7 streaming service analyses
├── ref/
│   ├── qbz-client.md                 # QBZ client deep-dive (databases, radio, blacklist)
│   ├── api_routes.md                 # MMS API endpoints
│   └── constants.md                  # Extracted constants
├── data/  (mms.db [SQLite WAL], artwork/)
├── server/                            # Go backend
│   ├── cmd/mms/main.go               # Entry point
│   ├── internal/api/  (router.go, handlers.go)
│   ├── internal/config/config.go
│   ├── internal/db/  (db.go, models.go, repository.go)
│   ├── internal/scanner/scanner.go / search/search.go
│   ├── internal/stream/stream.go      # Audio streaming + Range
│   └── go.mod / Dockerfile            # Multi-stage Alpine + FFmpeg
└── web/                               # React frontend
    ├── src/api/ (client.ts, hooks.ts, types.ts)
    ├── src/components/ (Layout, NowPlayingBar, Sidebar, AlbumCard, TrackList)
    ├── src/pages/ (Home, Search, Library, Album, Artist)
    ├── src/store/player.ts            # Zustand + Howler.js
    └── src/hooks/useKeyboardShortcuts.ts
```

## Conventions
- Go: chi router + handler struct, zerolog, context-aware repository
- React: client.ts -> hooks.ts -> component pattern
- Zustand wraps Howler.js for global player state
- FTS5 queries append `*` for prefix matching
- UPSERT pattern for idempotent metadata updates
- No authentication (local network assumption)

## Hard Assertions
```
NEVER: Commit QBZ auth credentials (~/.config/qbz/.qbz-auth)
NEVER: Block on audio streaming (write timeout disabled for Range requests)
ALWAYS: Use html5: true in Howler.js (streaming, not load-into-memory)
ALWAYS: Use WAL mode for SQLite
ALWAYS: Use parameterized SQL queries (no string interpolation)
ALWAYS: Record play history on track end
```

## Known Issues
- MMS: No authentication (local network only)
- MMS: CORS allows all origins (development mode)
- MMS: No test files exist yet
- MMS: Transcoding not yet implemented (cache path configured)
- MMS: Artist page partially complete, Settings page is stub
- QBZ: No track-level blacklist (artist-level only) — feature request candidate
