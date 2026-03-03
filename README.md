# Tunescape

A music taste analyzer that connects to your Spotify account and generates analytics and summaries based on your listening history.

## Prerequisites

Make sure you have the following installed before anything else:

- [Docker](https://www.docker.com/products/docker-desktop) and Docker Compose
- [Go 1.23+](https://go.dev/dl/) (only needed if running outside Docker)
- [Make](https://www.gnu.org/software/make/) (comes pre-installed on macOS and Linux, Windows users can use WSL)
- A Spotify Developer account with an app created at [developer.spotify.com](https://developer.spotify.com/dashboard)

## Getting Started

### 1. Clone the repository

```bash
git clone https://gitlab.com/Uranury/tunescape.git
cd tunescape
```

### 2. Set up your environment

Copy the example env file and fill in your values:

```bash
cp .env.example .env
```

Open `.env` and fill in the required values. See the [Environment Variables](#environment-variables) section below for details.

### 3. Start the project

```bash
make up
```

This starts the API, PostgreSQL, and Redis. Migrations run automatically on startup.

### 4. Verify everything is running

```bash
make logs
```

You should see `server starting` in the logs. If Postgres or Redis fail their healthchecks, wait a few seconds and run `make logs` again.

## Common Commands

| Command | Description |
|---|---|
| `make up` | Start all services in the background |
| `make down` | Stop all services |
| `make build` | Rebuild and start (use after pulling changes) |
| `make restart` | Restart only the API container |
| `make logs` | Tail API logs |

## Environment Variables

Copy `.env.example` to `.env` and fill in the following:

| Variable | Description | Required |
|---|---|---|
| `DB_USER` | Postgres username | Yes |
| `DB_PASSWORD` | Postgres password | Yes |
| `DB_NAME` | Postgres database name | Yes |
| `DB_HOST` | Postgres host (use `postgres` inside Docker) | Yes |
| `DB_PORT` | Postgres port | Default: `5432` |
| `DB_DRIVER` | Database driver | Use `postgres` |
| `DB_SSLMODE` | SSL mode | Default: `disable` |
| `REDIS_ADDR` | Redis address (use `redis:6379` inside Docker) | Yes |
| `LISTEN_ADDR` | Address the API listens on | Default: `:8080` |
| `JWT_KEY` | Secret key for signing JWTs, make it long and random | Yes |
| `MIGRATIONS_PATH` | Path to migration files | Default: `./migrations` |
| `ALLOWED_ORIGINS` | Comma-separated list of allowed CORS origins | Yes |

## Project Structure

```
.
├── cmd/
│   └── api/            # Application entrypoint (main.go)
├── internal/
│   ├── app/            # HTTP server setup and routing
│   ├── auth/           # JWT and refresh token logic
│   └── infra/          # Infrastructure wiring (DB, Redis, config)
├── pkg/
│   ├── config/         # Config loading
│   └── database/       # DB init, migrations, transaction support
├── migrations/         # SQL migration files
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## Gitflow

We follow a standard gitflow. Never push directly to `main` or `dev`.

- Branch off `dev` for your feature: `git checkout -b feature/your-feature-name`
- Open a PR into `dev` when ready for review
- `main` is only updated via releases from `dev`

Commit messages should be clear and descriptive. Prefix with the type of change: `feat:`, `fix:`, `chore:`, `docs:`.

## Troubleshooting

**Postgres or Redis not ready on first start** — Run `make down && make up`. On the first run Docker needs to pull images which can cause timing issues.

**Port 8080 already in use** — Change `LISTEN_ADDR` in your `.env` to another port like `:8081` and update the port mapping in `docker-compose.yml` accordingly.

**Migrations failing** — Make sure `DB_HOST` is set to `postgres` (the Docker service name) and not `localhost` when running inside Docker.