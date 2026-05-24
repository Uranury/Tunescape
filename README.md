# Tunescape

Tunescape is a Spotify music taste analyzer. Connect your Spotify account, capture snapshots of your listening history, and get analytics on your audio profile — energy, valence, danceability, and more. Track how your taste changes over time, compete on leaderboards, compare with friends, generate a PDF report, and create playlists directly from your top tracks.

## Features

- **Spotify OAuth** — connect and disconnect your Spotify account
- **Snapshots** — capture your current top tracks with audio features (energy, valence, danceability, acousticness, tempo, etc.)
- **Analytics** — averaged audio feature scores across your latest snapshot
- **Trends** — chart how your taste evolves across snapshots over time
- **Leaderboards** — global Redis-backed rankings for valence, energy, danceability, and acousticness
- **Friends** — send/accept friend requests, compare taste compatibility scores side by side
- **Playlists** — generate a Spotify playlist from your top snapshot tracks, embedded and playable in-app
- **PDF Report** — download a formatted report of your top tracks and leaderboard rankings
- **Background worker** — automatically re-captures snapshots every 24 hours per user

## Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop) and Docker Compose
- [Go 1.25+](https://go.dev/dl/) (only needed to run tests or build outside Docker)
- [Make](https://www.gnu.org/software/make/) (pre-installed on macOS/Linux; Windows users use WSL)
- A Spotify Developer app at [developer.spotify.com](https://developer.spotify.com/dashboard)

## Getting Started

### 1. Clone the repository

```bash
git clone https://gitlab.com/Uranury/tunescape.git
cd tunescape
```

### 2. Set up environment variables

```bash
cp .env.example .env
```

Open `.env` and fill in the required values (see [Environment Variables](#environment-variables)).

### 3. Start the project

```bash
make up
```

This starts the API, PostgreSQL, and Redis. Migrations run automatically on startup. The frontend is served at `http://localhost:8080`.

### 4. Verify everything is running

```bash
make logs
```

You should see `server starting` in the logs. If Postgres or Redis fail their health checks on first run, wait a few seconds and try again.

## Common Commands

| Command | Description |
|---|---|
| `make up` | Start all services in the background |
| `make down` | Stop all services |
| `make build` | Rebuild and restart (use after pulling changes) |
| `make restart` | Restart only the API container |
| `make logs` | Tail API logs |
| `make migrate` | Run migrations manually |
| `make fmt` | Format code with goimports |

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `DB_USER` | PostgreSQL username | — |
| `DB_PASSWORD` | PostgreSQL password | — |
| `DB_NAME` | PostgreSQL database name | — |
| `DB_HOST` | PostgreSQL host (`postgres` inside Docker) | — |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_DRIVER` | Database driver | `postgres` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `REDIS_ADDR` | Redis address (`redis:6379` inside Docker) | — |
| `LISTEN_ADDR` | Address the API listens on | `:8080` |
| `JWT_KEY` | Secret key for signing JWTs — make it long and random | — |
| `SPOTIFY_CLIENT_ID` | Spotify app client ID | — |
| `SPOTIFY_CLIENT_SECRET` | Spotify app client secret | — |
| `SPOTIFY_REDIRECT_URL` | OAuth callback URL registered in your Spotify app | — |
| `MIGRATIONS_PATH` | Path to migration files | `./migrations` |
| `ALLOWED_ORIGINS` | Comma-separated CORS origins | — |

## API Overview

All protected routes require a `Bearer` JWT in the `Authorization` header. Swagger docs are available at `/swagger/index.html`.

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/auth/signup` | No | Create an account |
| `POST` | `/auth/login` | No | Log in, receive JWT + refresh token |
| `POST` | `/auth/logout` | No | Invalidate refresh token |
| `POST` | `/auth/refresh` | No | Rotate JWT using refresh token |
| `GET` | `/auth/spotify/login` | No | Redirect to Spotify OAuth consent |
| `GET` | `/auth/spotify/callback` | No | OAuth callback, stores Spotify tokens |
| `GET` | `/me/profile` | Yes | Get current user profile |
| `DELETE` | `/me/spotify` | Yes | Disconnect Spotify account |
| `POST` | `/me/snapshots` | Yes | Capture a new snapshot from Spotify |
| `GET` | `/me/snapshots` | Yes | List all snapshots |
| `GET` | `/me/snapshots/:id` | Yes | Get a snapshot with its tracks |
| `GET` | `/me/trends` | Yes | Trend data across all snapshots |
| `GET` | `/me/report` | Yes | Download PDF report |
| `POST` | `/me/playlists/top-tracks` | Yes | Create a Spotify playlist from latest snapshot |
| `GET` | `/analytics/top-tracks` | Yes | Audio feature averages for latest snapshot |
| `GET` | `/leaderboards/:feature` | No | Global leaderboard (`valence`, `energy`, `danceability`, `acousticness`) |
| `GET` | `/users/lookup` | Yes | Look up a user by display name |
| `POST` | `/friends/requests` | Yes | Send a friend request |
| `GET` | `/friends/requests` | Yes | List incoming friend requests |
| `POST` | `/friends/requests/:id/accept` | Yes | Accept a friend request |
| `POST` | `/friends/requests/:id/reject` | Yes | Reject a friend request |
| `GET` | `/friends` | Yes | List friends |
| `GET` | `/friends/:friend_id/compare` | Yes | Compare taste scores with a friend |
| `GET` | `/friends/:friend_id/playlists` | Yes | View a friend's playlists |

## Project Structure

```
.
├── cmd/api/                  # Entrypoint (main.go)
├── internal/
│   ├── app/                  # HTTP server, routing, graceful shutdown
│   ├── auth/                 # JWT issuance, refresh token rotation
│   ├── spotify/              # OAuth2 flow, Spotify API client, token storage
│   ├── snapshot/             # Capture and store listening history
│   ├── analytics/            # Audio feature scoring
│   ├── trends/               # Trend analysis across snapshots
│   ├── leaderboard/          # Redis sorted-set leaderboards
│   ├── playlist/             # Playlist creation via Spotify API
│   ├── friends/              # Friend requests, taste comparison
│   ├── report/               # PDF report generation (gofpdf)
│   ├── worker/               # Background snapshot worker (24h tick)
│   ├── middleware/           # JWT auth, Redis-backed rate limiting
│   ├── cache/                # Redis cache wrapper
│   ├── user/                 # User profile
│   ├── track/                # Track model and repository
│   ├── reccobeats/           # Reccobeats API client (recommendations)
│   └── infra/                # Shared deps (DB, Redis, logger, config)
├── pkg/
│   ├── config/               # Config loading via cleanenv
│   ├── database/             # Connection, migrations, TxProvider
│   ├── apperrors/            # Sentinel errors
│   └── validation/           # Input validation helpers
├── frontend/
│   ├── index.html
│   ├── css/main.css
│   └── js/                   # Vanilla JS modules (auth, dashboard, analytics, playlist, friends, leaderboard)
├── migrations/               # SQL migrations (9 versions, run automatically)
├── docs/                     # Swagger-generated API docs
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## Architecture

Every feature follows the same three-layer pattern:

```
Handler → Service → Repository
```

- **Handler** — validates input, calls service, writes JSON response (Gin)
- **Service** — business logic, orchestrates repositories and external clients
- **Repository** — SQL queries via `sqlx`, receives a `TxProvider` for transactions

`internal/infra/deps.go` holds shared infrastructure (`Deps`). All wiring is done manually in `cmd/api/main.go` — no DI framework.

### Authentication Flow

1. User registers or logs in → receives a short-lived JWT and a refresh token stored in PostgreSQL
2. JWT is passed as `Authorization: Bearer <token>` on protected routes
3. `POST /auth/refresh` rotates the refresh token and issues a new JWT
4. Separately, `/auth/spotify/login` initiates Spotify OAuth; the callback stores Spotify access/refresh tokens in the DB and updates the user's profile

### Background Worker

`internal/worker/snapshot_worker.go` runs as a goroutine started at boot. Every 24 hours it iterates over all users with connected Spotify accounts, fetches their top tracks, and persists a new snapshot. Includes panic recovery to keep the worker alive through unexpected errors.

### Rate Limiting

`GET /leaderboards/:feature` is rate-limited to 60 requests per minute per IP, enforced with a Redis sliding window.

## Running Tests

Tests require a local PostgreSQL 16 instance and Redis 7:

- PostgreSQL: `localhost:5432`, database `tunescape_test`, user `test`, password `test`
- Redis: `localhost:6379`

```bash
go test ./... -race -count=1
```

Run a single test:

```bash
go test ./internal/<module>/... -run TestName
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for branching conventions, commit prefixes, and code style. Never push directly to `main` or `dev` — branch off `dev` as `feature/name` or `fix/name` and open a PR.

## Troubleshooting

**Postgres or Redis not ready on first start** — Run `make down && make up`. On first run Docker needs to pull images which can cause timing issues.

**Port 8080 already in use** — Change `LISTEN_ADDR` in `.env` to another port (e.g. `:8081`) and update the port mapping in `docker-compose.yml`.

**Migrations failing** — Make sure `DB_HOST` is set to `postgres` (the Docker service name) when running inside Docker, not `localhost`.

**Spotify embed shows grayed out** — This is a propagation delay on Spotify's end immediately after playlist creation. The embed loads automatically after a short delay.
