# AGENTS.md

## Run
**Do NOT run the server** — the user runs it manually.
```
go run cmd/api/main.go
```
Requires `.env` with `PORT`, `DATABASE_URL`, `APP_ENV`.

## Architecture
- `cmd/api/main.go` — entry point
- `internal/bootstrap/app.go` — wires Store → Service → Controller → Routes
- `internal/routes/` — all routes mounted under `/api/v1`
- `internal/middlewares/auth_middleware.go` — JWT auth via cookie (`access_token`); `JWT_SECRET` is read via `os.Getenv` (not via `pkg/env`)

## Auth
`/api/v1/auth/login` (POST) returns the cookie. Protected routes use `middlewares.Protected()`.

## DB
SQL schema at `db/vehicle-control.sql`. Run migrations manually against Neon PostgreSQL.

## Notes
- `pkg/env/env.go` does not load `JWT_SECRET` into the `Config` struct — middleware reads it directly from env.
- Vehicle routes are referenced in `routes.go` but not yet implemented.
- No test suite, no linting config, no CI.