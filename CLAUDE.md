# Mark's Music Solutions — Project Guide

Self-hosted music streaming server + ethical streaming research. Go backend (chi, SQLite FTS5, zerolog), React frontend (Zustand, Howler.js, Tailwind). Docker-ready, fully functional.

## Quick Reference
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
NEVER: Block on audio streaming (write timeout disabled for Range requests)
ALWAYS: Use html5: true in Howler.js (streaming, not load-into-memory)
ALWAYS: Use WAL mode for SQLite
ALWAYS: Use parameterized SQL queries (no string interpolation)
ALWAYS: Record play history on track end
```
## Known Issues
- No authentication (local network only)
- CORS allows all origins (development mode)
- No test files exist yet
- Transcoding not yet implemented (cache path configured)
- Artist page partially complete, Settings page is stub

See ref/ for API routes, extracted constants, and lookup tables.
