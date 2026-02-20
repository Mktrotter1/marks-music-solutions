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
