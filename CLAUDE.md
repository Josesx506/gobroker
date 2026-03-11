# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

All commands run from the `backend/` directory.

```bash
# Run the application
go run main.go

# Build binary
go build -o tmp/main

# Run with hot reload (requires air, pre-installed in devcontainer)
air

# Database migrations (auto-run on startup, but can run manually)
goose -dir migrations postgres "$DATABASE_URL" up
goose -dir migrations postgres "$DATABASE_URL" down
```

**No Makefile** — there are no make targets.

## Dev Environment

The project uses a devcontainer (`.devcontainer/`) with:
- PostgreSQL 15 (via docker-compose)
- `air` for hot reload
- `goose` for migrations

Copy `.devcontainer/.dev.env` values into `backend/.env` when running locally outside devcontainer.

## Architecture

GoBroker is a real-time event streaming backend. When a temperature reading is inserted into PostgreSQL, a trigger fires `pg_notify`, which a Go listener picks up and broadcasts to connected SSE clients filtered by `location_id`.

**Data flow:**
```
INSERT → PostgreSQL TRIGGER → pg_notify('temp_events') → Listener goroutine → Broker → SSE clients
```

### Key packages (`backend/internal/`)

- **`app/`** — Wires dependencies together: opens DB, runs embedded migrations, creates logger
- **`store/`** — Two responsibilities:
  - `database.go`: Connection pool setup + Goose migrations via embedded FS
  - `listener.go`: Long-running goroutine with exponential backoff (1s–60s) that LISTENs on `temp_events` and calls `broker.Broadcast()`
- **`broker/`** — In-memory pub/sub. Routes messages only to clients subscribed to the matching `location_id`
- **`api/`** — chi/v5 router setup; currently minimal, ready for endpoint expansion

### Database schema

- `public.temperature_readings` (`id`, `location_id`, `device_id`, `value`, `created_at`)
- `rtmfuncs.notify_temp_insert()` trigger function sends JSON via `pg_notify` after each insert
- Migrations are embedded in the binary (`backend/migrations/`) and auto-applied on startup

### Dependencies

- **HTTP**: `go-chi/chi/v5`
- **PostgreSQL**: `jackc/pgx/v4`
- **Migrations**: `pressly/goose/v3`
- **Config**: `joho/godotenv`
