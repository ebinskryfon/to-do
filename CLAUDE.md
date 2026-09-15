# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A teaching-oriented Go clean-architecture Todo API (module `todo`). It is being built incrementally to
demonstrate strict layering to a junior developer, mirroring the (code-free, Markdown-only) architecture
spec from `/home/ebinb/employee-360`. **`ARCHITECTURE.md` at the repo root is the design doc and build
plan — read it before making structural changes.** It describes the intended package layout, build order,
and a suggested reading order through the codebase. Treat it as the source of intent, but note the
implementation has already diverged from it in a few places (see "Known divergences from ARCHITECTURE.md"
below) — when in doubt, prefer what the code actually does.

Only the `Create` todo operation is fully wired end-to-end today (usecase interface + impl, handler route,
container wiring). `Get`/`List`/`Update`/`Delete`/`Toggle` exist only as methods on the `domain/repository`
port; their `usecase/interfaces` and `usecase/todo` counterparts, handler methods, and routes have not been
added yet. `cmd/bootstrap` is a stub directory (`.gitkeep` only, no `main.go`).

## Commands

```bash
make run              # go run ./cmd/api
make build            # go build -o bin/api ./cmd/api
make test             # go test ./...
make tidy             # go mod tidy
make migrate-up       # go run ./cmd/migrate up   (also creates the DB if missing)
make migrate-down     # go run ./cmd/migrate down
make bootstrap        # go run ./cmd/bootstrap    (not implemented yet)
make generate         # swag init -g cmd/api/main.go -o docs  (regenerate docs/ from handler annotations)
```

Run a single test:
```bash
go test ./internal/usecase/todo/ -run TestCreateTodoUsecase_Success -v
```

`cmd/migrate/main.go` is a standalone CLI built directly on the `golang-migrate/migrate/v4` library (no
external `migrate` binary needed). It supports more than `up`/`down`: `status`, `force <version>`, and
`create <name>` (auto-numbers the next migration prefix by scanning `migrations/`). `migrate up` also
creates the target database first if it doesn't exist (connects to the `postgres`/`template1` admin DB to
issue `CREATE DATABASE`). DB URL resolution order: `MIGRATIONS_DATABASE_URL` env → `DATABASE_URL` env →
`config.Database.DSN()`.

Local setup: `cp .env.example .env`, adjust DB credentials, then `make migrate-up && make run` (server
defaults to port 8080).

## Architecture

**Dependency rule, enforced one-directionally:** Infrastructure → Delivery → Usecase → Domain. Domain has
zero framework imports. Usecases import only `domain/entity` + `domain/repository`, never gorm/gin.
Persistence implements the domain repository port and is the only place `gorm` appears outside
`infrastructure/database`. Handlers import usecase interfaces + `response`, never gorm/persistence
directly. `infrastructure/container` is the one file allowed to import concrete types from every layer to
wire them together bottom-up (repo → usecases → handlers).

**Usecase ports/impl split** (same pattern as `domain/repository` vs `infrastructure/persistence`, one
layer up): `internal/usecase/interfaces/` holds one file per operation containing only an interface
(`CreateTodoUsecase`, etc.), importing just `domain/entity` + `context`. `internal/usecase/todo/` holds the
matching implementation, one file per operation, with an unexported struct whose constructor
(`NewCreateTodoUsecase`) returns the **interface type**, not the concrete struct — so handlers and the
container only ever name the interface, and tests can substitute a mock repository with zero DB/HTTP
involvement (see `internal/usecase/todo/create_todo_test.go`). Follow this same interface+impl split when
adding any new usecase — do not add bare structs directly under `usecase/todo` without a matching port in
`usecase/interfaces`.

**Entities carry audit + soft-delete fields** (`internal/domain/entity/base.go`): `BaseModel` (id,
timestamps, `DeletedAt`) → `AuditableEntity` (adds `CreatedBy`/`UpdatedBy`) → `SoftDeletableEntity` (adds
`IsActive`/`DeletedBy`). `Todo` embeds `SoftDeletableEntity`. Deletes are soft (`is_active=false`,
`deleted_at` set); `GetByID`/`List`/`Update` all filter on `is_active = true AND deleted_at IS NULL` at the
persistence layer (`internal/infrastructure/persistence/todo_repo.go`). The persistence layer uses an
internal `todoModel` (gorm tags) with `toEntity()`/`fromEntity()` conversions, keeping `entity.Todo` itself
framework-free. `gorm.ErrRecordNotFound` and zero-`RowsAffected` updates/deletes both translate to
`domainerrors.ErrTodoNotFound` at this boundary.

**Config** (`pkg/config/config.go`): Viper-based, precedence real env vars > `.env` (loaded via
`godotenv.Load()`, tried at both `.env` and `../.env`) > `pkg/config/config.yaml` > defaults. `Config`
holds `AppEnv`, `ServerPort`, and a nested `Database DatabaseConfig` (host/port/user/password/name/sslmode
plus connection-pool settings: `MaxOpenConns`, `MaxIdleConns`, `ConnMaxIdleTime`, `ConnMaxLifetime`).
`DatabaseConfig.DSN()` returns `Database.URL` verbatim if set, otherwise builds a `postgres://` URL from
the individual fields. Flat `DB*` fields on `Config` are kept as aliases for backwards compatibility —
prefer `Config.Database` in new code.

**Response envelope** (`internal/delivery/http/response/response.go`): every JSON response has the shape
`{success, data?, error?{code, message}, meta?{total}}`. Handlers must go through
`Success`/`Created`/`Paginated`/`Error` — never call `c.JSON` directly. `Error` calls
`AbortWithStatusJSON`, so no downstream handler/middleware writes to the response afterward.

**Error → HTTP mapping** happens only in handlers (`handleError` in `todo_handler.go`), via `errors.Is`
against `domainerrors.ErrInvalidInput` (400) / `domainerrors.ErrTodoNotFound` (404), default 500. Usecase
and repository layers never reference HTTP status codes.

**Middleware chain**, registered in `internal/delivery/http/routes.go` in this order: `RequestID` → `Logger`
→ `CORS` → `Recovery` → `AuditContext`. `AuditContext` (`internal/delivery/http/middleware/audit_context.go`)
stashes client IP/User-Agent into both the Gin context and a `pkg/audit.Context` on the request's
`context.Context`, for future use by usecases that need to stamp `CreatedBy`/`UpdatedBy`.

**Server lifecycle** (`internal/infrastructure/server/server.go`): `Run()` starts `ListenAndServe` in a
goroutine, blocks on `signal.NotifyContext` for SIGINT/SIGTERM, then calls `Shutdown` with a 10s timeout.

**Swagger**: handler methods carry `swaggo` annotations; `make generate` regenerates `docs/` from them.
`cmd/api/main.go` carries the top-level API annotations (title/version/host/etc). Docs are served at
`/swagger/*any`.

## Known divergences from ARCHITECTURE.md

- ARCHITECTURE.md's original plan explicitly excluded audit logging ("no audit logging" — see its Context
  section) for this simplified teaching app, but the current code has audit/soft-delete fields on every
  entity and an `AuditContext` middleware + `pkg/audit` package. Match the existing audit/soft-delete
  pattern for any new entity or migration, not the original no-audit plan.
- `cmd/migrate/main.go` is considerably more featureful than the plan's `up`/`down`-only sketch (adds
  `status`, `force`, `create`, and auto-creates the target database).
- `pkg/config` uses a nested `DatabaseConfig` struct with connection-pool settings, not the flat
  `Config{DBHost, DBPort, ...}` struct the plan sketched (flat fields are kept only as aliases).
- Migration file naming is `NNN_name.up.sql`/`.down.sql` (3-digit, underscore) rather than the plan's
  6-digit `NNNNNN_name` — `cmd/migrate create` follows the 3-digit convention actually in use.
