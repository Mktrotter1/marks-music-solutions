# QBZ — Native Qobuz Client for Linux

Installed 2026-03-02 from AUR (`qbz-bin 1.1.18-1`). Built with Tauri (Rust + WebKit).
Upstream: https://github.com/vicrodh/qbz

## Installation

| Field | Value |
|-------|-------|
| Package | `qbz-bin` (AUR) |
| Binary | `/usr/bin/qbz` (68 MiB) |
| Desktop Entry | `/usr/share/applications/qbz.desktop` |
| WM Class | `com.blitzfc.qbz` |
| Audio | PipeWire (pipewire-alsa, pipewire-pulse) |
| Dependencies | webkit2gtk-4.1, gtk3, alsa-lib, libappindicator-gtk3, libxkbcommon, openssl |

## Data Locations

| Purpose | Path |
|---------|------|
| Auth credentials | `~/.config/qbz/.qbz-auth` |
| App data root | `~/.local/share/qbz/` |
| User data (ID: 10620517) | `~/.local/share/qbz/users/10620517/` |
| Tauri/WebKit data | `~/.local/share/com.blitzfc.qbz/` |
| Cache (artwork, playback) | `~/.cache/qbz/` |

## SQLite Databases

All state is in SQLite. Key databases:

### Radio Engine (`~/.local/share/qbz/radio/radio_engine.db`)

| Table | Purpose |
|-------|---------|
| `radio_session` | Active radio stations — seed type (artist/track), spacing, reseed interval |
| `radio_pool` | Candidate tracks per session — `used` flag prevents replay within session |
| `radio_history` | Played tracks log per session |

Radio config per session: `artist_spacing=5` (min tracks between same artist), `reseed_every=25` (fetch new candidates after 25 plays).

### Artist Blacklist (`~/.local/share/qbz/users/10620517/artist_blacklist.db`)

| Table | Purpose |
|-------|---------|
| `artist_blacklist` | Blocked artists (artist_id, artist_name, added_at, notes) |
| `blacklist_settings` | Feature toggle (enabled=1 by default) |

Blacklist is checked in radio, queue, recommendations, and discovery paths. Log output: `[V2/Blacklist] Filtered`.

**No track-level blacklist exists** — only artist-level.

### Tauri IPC Commands (relevant)

| Command | Purpose |
|---------|---------|
| `v2_add_to_artist_blacklist` | Block artist from all radio/reco |
| `v2_remove_from_artist_blacklist` | Unblock artist |
| `v2_get_artist_blacklist` | List blocked artists |
| `v2_set_blacklist_enabled` | Toggle blacklist on/off |
| `v2_clear_artist_blacklist` | Remove all entries |
| `v2_create_artist_radio` | Start radio seeded from artist |
| `v2_create_track_radio` | Start radio seeded from track |
| `v2_create_album_radio` | Start radio seeded from album |

### Other User Databases

| Database | Purpose |
|----------|---------|
| `library.db` | Local music library (playlists, folders, local tracks, custom covers) |
| `favorites_cache.db` | Favorited albums, artists, tracks |
| `session.db` | Current player state + queue |
| `playback_preferences.db` | Autoplay mode, context icon, session persistence |
| `audio_settings.db` | Output device, exclusive mode, DAC passthrough, gapless, normalization |
| `cache/api_cache.db` | Qobuz API response cache |
| `cache/musicbrainz_cache.db` | MusicBrainz metadata enrichment |
| `cache/listenbrainz_cache.db` | ListenBrainz integration cache |
| `reco/events.db` | Recommendation engine (play events, scores, track/album/artist metadata) |

### Recommendation Engine (`reco/events.db`)

- Event types: `play`, `playlist_add`
- Score types: `all`
- Tracks metadata including MusicBrainz IDs (ISRC, mbid)

## Integrations

- **Last.fm**: Scrobbling support (auth + session key storage)
- **ListenBrainz**: Listen submission support
- **MusicBrainz**: Artist/track metadata resolution
- **Discogs**: Album artwork search
- **Chromecast**: Cast playback
- **DLNA**: Network playback
- **AirPlay**: Apple device playback
- **Plex**: Plex library integration

## How Radio Works

1. User creates a radio station seeded from an **artist**, **track**, or **album**
2. QBZ fetches candidate tracks into `radio_pool` with sources (`seed_tracks`, `similar_artist`) and `distance` scores
3. Tracks are selected respecting `artist_spacing` (avoid same artist too close) and randomized via `rng_seed`
4. Played tracks marked `used=1` in pool — won't replay in that session
5. After `reseed_every` (25) plays, new candidates are fetched
6. Artist blacklist is applied as a filter across all candidate selection

## Network & DNS (Tailscale Coexistence)

**Problem (fixed 2026-03-25):** QBZ would intermittently fail with "bad gateway" errors. Root cause: Tailscale's MagicDNS was the sole system DNS resolver (`/etc/resolv.conf` → `100.100.100.100`). Any MagicDNS hiccup meant Qobuz API/CDN domains couldn't resolve.

**Fix:** Enabled `systemd-resolved` as the system DNS stub resolver. Tailscale auto-detects resolved via D-Bus and registers itself only for tailnet domains (`*.ts.net`, `100.x` PTR). Internet DNS (including Qobuz) routes through normal upstreams with fallback.

| Config file | Purpose |
|-------------|---------|
| `/etc/systemd/resolved.conf.d/dns-fallback.conf` | Cloudflare (1.1.1.1) primary, Google (8.8.8.8) fallback, DNS-over-TLS opportunistic |
| `/etc/resolv.conf` | Symlink → `/run/systemd/resolve/stub-resolv.conf` (127.0.0.53) |

DNS routing:

| Query type | Resolver | Interface |
|------------|----------|-----------|
| `*.ts.net`, Tailscale IPs | MagicDNS (100.100.100.100) | tailscale0 |
| Everything else (Qobuz, Last.fm, etc.) | Router DNS / Cloudflare / Google | enp7s0 or wlan0 |

QBZ connects to Qobuz via Akamai CDN (streaming, API, artwork) and Cloudflare (scrobbling integrations). All connections go out on the local interface, not through Tailscale.

**If bad gateway returns:** Run `resolvectl status` — verify `tailscale0` shows `Default Route: no` and your LAN interface shows `Default Route: yes`. If Tailscale has taken over again, restart resolved: `sudo systemctl restart systemd-resolved && sudo systemctl restart tailscaled`.

## Removing Songs/Artists from Radio

| Goal | Method |
|------|--------|
| Block an **artist** permanently | Use artist blacklist (UI: right-click artist or Settings) |
| Skip a **track** in current session | Skip it — marked `used` in pool, won't repeat this session |
| Block a **track** permanently | Not supported — feature request needed on GitHub |
