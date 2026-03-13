# TODO — Mark's Music Solutions

## High Priority

- [ ] Add authentication layer (currently no auth — local network assumption)
- [ ] Restrict CORS origins (currently `Access-Control-Allow-Origin: *`)
- [ ] Implement transcoding via FFmpeg (cache path configured in `config.yaml` but `stream.go` only serves raw files)

## In Progress

- [ ] Complete Artist page (displays albums but missing top tracks, bio, related artists)
- [ ] Build Settings page (Sidebar links to `/settings` but no route or page component exists)

## Backlog

- [ ] Add server-side tests (`server/internal/` has zero `_test.go` files)
- [ ] Add frontend tests (no test setup in `web/`)
- [ ] Implement play history UI (API endpoint `POST /tracks/{id}/play` exists, no frontend display)
- [ ] Add playlist management UI (API endpoints exist: CRUD playlists, add tracks)
