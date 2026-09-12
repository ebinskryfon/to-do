# Plan: Simple Go Clean-Architecture Todo API (teaching example, mirrors employee-360)

## Context

The goal is a teaching artifact for a junior developer learning Go clean architecture. The user asked to base it on the `employee-360` backend at `/home/ebinb/employee-360`. That repo turned out to be **planning-only — no Go code exists there at all**; it's a set of Markdown specs (`plan/architecture/backend.md`, `.agents/rules/clean-architecture.md`, `.agents/skills/new-backend-feature/SKILL.md`) written for an AI agent to scaffold from later. Those specs are precise about folder layout, layering rules, DI pattern, response envelope, middleware chain, and migration conventions, so this plan mirrors that spec (not literal code, since none exists) and simplifies away everything that's overkill for a single-entity todo app: no multi-tenancy, no JWT/RBAC/auth, no audit logging. It keeps `cmd/migrate` and `cmd/bootstrap` as standalone binaries (scaled down: bootstrap seeds sample todos, not tenants/roles/admin).

Decisions already confirmed with the user:
- Go module name: simple (`todo`), not a `github.com/org/...` path
- Single `Todo` entity, full CRUD + toggle-done (no second related entity)
- Project lives in `/home/ebinb/to-do` (currently empty, not a git repo)
- No Docker Compose — connects to an already-running local Postgres via `.env`/Viper config

Stack (mirrors employee-360's stated tech choices, minus auth/JWT/RBAC/audit/mail which don't apply here): Gin, GORM + Postgres driver, golang-migrate, Viper, go-playground/validator, zerolog, google/uuid, testify.

Dependency rule enforced throughout: **Infrastructure → Delivery → Usecase → Domain**. Domain has zero framework imports. Usecase imports only domain (interfaces + entities), never gorm/gin. Persistence implements the domain repository port and is the only place gorm appears outside `database/postgres.go`. Handlers import usecase + response, never gorm/persistence directly. The DI container is the single file allowed to import concrete types from every layer, to wire them together.

## Project structure

```
/home/ebinb/to-do/
├── cmd/
│   ├── api/main.go
│   ├── migrate/main.go           # migration runner CLI (golang-migrate as a library)
│   └── bootstrap/main.go         # idempotent seeder CLI (sample todos for local dev)
├── internal/
│   ├── domain/
│   │   ├── entity/todo.go
│   │   ├── repository/todo_repository.go
│   │   └── errors/errors.go
│   ├── usecase/
│   │   ├── interfaces/            # ports: one interface per operation, zero implementation
│   │   │   ├── create_todo.go
│   │   │   ├── get_todo.go
│   │   │   ├── list_todos.go
│   │   │   ├── update_todo.go
│   │   │   ├── delete_todo.go
│   │   │   └── toggle_todo.go
│   │   └── todo/                  # implementation of each interface above
│   │       ├── create_todo.go
│   │       ├── get_todo.go
│   │       ├── list_todos.go
│   │       ├── update_todo.go
│   │       ├── delete_todo.go
│   │       └── toggle_todo.go
│   ├── infrastructure/
│   │   ├── database/postgres.go
│   │   ├── persistence/todo_repo.go
│   │   ├── container/container.go
│   │   └── server/server.go        # HTTP server lifecycle: start + graceful shutdown
│   └── delivery/http/
│       ├── handlers/todo_handler.go
│       ├── middleware/{request_id.go,logger.go,cors.go,recovery.go}
│       ├── response/response.go
│       └── routes.go
├── pkg/
│   ├── config/
│   │   ├── config.go
│   │   └── config.yaml.example
│   ├── logger/logger.go
│   └── validator/validator.go
├── migrations/
│   ├── 000001_create_todos.up.sql
│   └── 000001_create_todos.down.sql
├── .env.example
├── Makefile
├── go.mod
└── go.sum
```

This mirrors the same pattern already used for the repository layer (`domain/repository` = port, `infrastructure/persistence` = implementation), applied one level up: `usecase/interfaces` = port, `usecase/todo` = implementation. It also restores `cmd/migrate`, `cmd/bootstrap`, and `internal/infrastructure/server/` from employee-360's spec instead of assuming an externally-installed `migrate` CLI binary, skipping seeding entirely, or handling server start/shutdown inline in `main.go`.

## Build order

1. **Skeleton**: `go mod init todo`, create directory tree. (Go version is whatever `go mod tidy` pins based on dependency requirements — currently 1.25, since `gin-gonic/gin` v1.12 requires it.)
2. **Domain** (zero deps except `uuid`/`time`/`errors`):
   - `entity/todo.go`: `Todo{ID uuid.UUID, Title, Description string, Completed bool, CreatedAt, UpdatedAt time.Time}` — plain struct, no gorm tags.
   - `repository/todo_repository.go`: `TodoRepository` interface — `Create`, `GetByID`, `List`, `Update`, `Delete`, every method takes `ctx context.Context` first. No separate `Toggle` method — usecase composes `GetByID`+`Update`.
   - `errors/errors.go`: sentinel errors `ErrTodoNotFound`, `ErrInvalidInput`.
3. **Config**: `pkg/config/config.go` (Viper: `Config{AppEnv, ServerPort, DBHost, DBPort, DBUser, DBPassword, DBName, DBSSLMode}`, `Load() (*Config, error)`) — lives under `pkg/` alongside `pkg/logger` and `pkg/validator` since, like them, it's a small reusable utility with no dependency on any other internal layer, rather than an app-specific concern. `Load` calls `godotenv.Load()` first (ignoring a missing `.env` — it's optional) to pull `.env` into the real process environment, then `viper.AutomaticEnv()` picks those up automatically since Viper upper-cases lookup keys (`server_port` → `SERVER_PORT`), matching `.env`'s naming exactly. Precedence: real env vars > `.env` > `config.yaml` > defaults. `pkg/config/config.yaml.example`, top-level `.env.example`. Built before `database.go` since it's a compile dependency.
4. **Usecase**, split into a ports package and an implementation package — same pattern as `domain/repository` vs `infrastructure/persistence`, one level up:
   - `internal/usecase/interfaces/` — one file per operation, **interface only**, imports only `domain/entity` + `context`:
     ```go
     // internal/usecase/interfaces/create_todo.go
     type CreateTodoUsecase interface {
         Execute(ctx context.Context, title, description string) (*entity.Todo, error)
     }
     ```
     Files: `create_todo.go`, `get_todo.go`, `list_todos.go`, `update_todo.go`, `delete_todo.go`, `toggle_todo.go` — each declares exactly one interface. Return types vary: `(*entity.Todo, error)` for create/get/update/toggle, `([]entity.Todo, error)` for list, `error` for delete.
   - `internal/usecase/todo/` — one file per operation, **implementation only**, imports `domain/entity`, `domain/repository`, and the matching interface from `usecase/interfaces`:
     ```go
     // internal/usecase/todo/create_todo.go
     type createTodoUsecase struct { repo repository.TodoRepository }
     func NewCreateTodoUsecase(repo repository.TodoRepository) interfaces.CreateTodoUsecase {
         return &createTodoUsecase{repo: repo}
     }
     func (uc *createTodoUsecase) Execute(ctx context.Context, title, description string) (*entity.Todo, error) { ... }
     ```
     The struct is unexported; the constructor's return type is the **interface type from `usecase/interfaces`**, not the concrete struct — so nothing outside this package ever names `createTodoUsecase` directly. This is what lets handlers and the container depend on the interface and swap in a mock without touching this package.
5. **Infrastructure**:
   - `database/postgres.go`: `NewPostgresConnection(cfg config.Config) (*gorm.DB, error)` — builds DSN, opens via `gorm.Open`, pings immediately, returns error (fail-fast, no panic) if unreachable. No AutoMigrate — schema only from `migrations/`.
   - `persistence/todo_repo.go`: `todoRepository{db *gorm.DB}` implementing `repository.TodoRepository`. Uses `db.WithContext(ctx)` on every call. Translates `gorm.ErrRecordNotFound` → `domainerrors.ErrTodoNotFound` at this boundary. Use an internal `todoModel` struct with gorm/column tags + `toEntity()`/`fromEntity()` conversions, keeping `entity.Todo` itself framework-free — this makes the domain/infra boundary concrete for the junior dev to see. Add compile-time check `var _ repository.TodoRepository = (*todoRepository)(nil)`.
   - `container/container.go`: `NewContainer(db *gorm.DB) *Container` — manually wires repo → 6 usecases (each constructor from `usecase/todo` returns its interface type from `usecase/interfaces`, e.g. `interfaces.CreateTodoUsecase`) → `TodoHandler`, bottom-up. This is the only file importing persistence + usecase implementation + handlers together.
6. **`pkg/`**: `logger/logger.go` (`New(env string) zerolog.Logger`, pretty console in dev / JSON in prod), `validator/validator.go` (`Validate(s any) error` wrapping go-playground/validator).
7. **Delivery**:
   - `response/response.go`: `Envelope{Success bool, Data any, Error *ErrorBody, Meta *Meta}`; helpers `Success(c, data)` (200), `Created(c, data)` (201), `Paginated(c, data, total)` (200 + Meta.Total), `Error(c, status, code, message)` (writes envelope + `AbortWithStatusJSON`). Handlers never call `c.JSON` directly.
   - `middleware/`: `request_id.go` (generate/read `X-Request-ID`, first in chain), `logger.go` (zerolog structured request log), `cors.go` (permissive local-dev CORS), `recovery.go` (recover panics → `response.Error(500, "INTERNAL_ERROR", ...)`, last in chain). Registration order: **request_id → logger → cors → recovery**.
   - `handlers/todo_handler.go`: `TodoHandler` holds the 6 usecases **as their interface types from `usecase/interfaces`** (e.g. `createUC interfaces.CreateTodoUsecase`, not the concrete struct from `usecase/todo`) so handler tests can inject mocks without touching the repository; methods `Create/List/Get/Update/Delete/Toggle(c *gin.Context)`. Small request DTOs (`createTodoRequest`, `updateTodoRequest`) with validator tags in the same file. Each handler: bind JSON → validate → parse `:id` via `uuid.Parse` (400 on failure) → call usecase → map result/error via `response.*`. Never imports gorm/persistence.
   - `routes.go`: `SetupRouter(c *container.Container) *gin.Engine` — registers middleware, then:
     - `POST /api/v1/todos` → Create
     - `GET /api/v1/todos` → List
     - `GET /api/v1/todos/:id` → Get
     - `PUT /api/v1/todos/:id` → Update
     - `DELETE /api/v1/todos/:id` → Delete
     - `PATCH /api/v1/todos/:id/toggle` → Toggle
     - `GET /health` (outside `/api/v1`) → liveness check
8. **`internal/infrastructure/server/server.go`** + **`cmd/api/main.go`**: `server.go` wraps `*http.Server` with `New(router, port, log) *Server` and `(*Server) Run() error` — `Run` starts `ListenAndServe` in a goroutine, then blocks on `signal.NotifyContext(ctx, SIGINT, SIGTERM)` and, once a signal (or a listen error) arrives, calls `Shutdown(ctx)` with a 10s timeout so in-flight requests finish before exit. `main.go`'s boot sequence becomes: `config.Load()` → `logger.New(cfg.AppEnv)` → `database.NewPostgresConnection(*cfg)` (fatal on error) → `container.NewContainer(db)` → `http.SetupRouter(c)` → `server.New(router, cfg.ServerPort, log).Run()`.
9. **Migrations**: `migrations/000001_create_todos.up.sql` (`CREATE EXTENSION IF NOT EXISTS "pgcrypto"; CREATE TABLE todos (id uuid PRIMARY KEY DEFAULT gen_random_uuid(), title text NOT NULL, description text NOT NULL DEFAULT '', completed boolean NOT NULL DEFAULT false, created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now())`), `.down.sql` (`DROP TABLE IF EXISTS todos CASCADE`).
10. **`cmd/migrate/main.go`**: migration runner CLI, mirrors employee-360's `cmd/migrate`. Built as a thin wrapper around the `golang-migrate/migrate/v4` **library** (`database/postgres` + `source/file` drivers), not the external `migrate` CLI binary — so the junior dev doesn't need to install anything beyond Go/`go.mod`. Reads `config.Load()` for the DSN, takes a subcommand as `os.Args[1]`:
    ```go
    func main() {
        cfg, _ := config.Load()
        m, _ := migrate.New("file://migrations", buildDSN(cfg))
        switch os.Args[1] {
        case "up":   m.Up()
        case "down": m.Steps(-1)
        }
    }
    ```
11. **`cmd/bootstrap/main.go`**: idempotent seeder CLI, mirrors employee-360's `cmd/bootstrap` (there it seeds tenants/roles/admin; here, with no auth/tenant, it seeds a handful of sample todos for local dev convenience). Wires the same `container.NewContainer(db)` used by `cmd/api`, checks `listUC.Execute(ctx)` first and only inserts sample rows if the table is empty — matching the "seeder must be safe to re-run" rule from the spec (`seeder_test.go` in employee-360 verifies exactly this idempotency).
12. **Makefile**: `run`, `build`, `migrate-up`, `migrate-down`, `migrate-create`, `bootstrap`, `test`, `tidy` targets — `migrate-*` targets call `go run ./cmd/migrate <subcommand>` (the in-repo runner from step 10), not an externally-installed `migrate` binary:
    ```makefile
    run:            ; go run ./cmd/api
    build:          ; go build -o bin/api ./cmd/api
    migrate-up:     ; go run ./cmd/migrate up
    migrate-down:   ; go run ./cmd/migrate down
    bootstrap:      ; go run ./cmd/bootstrap
    test:           ; go test ./...
    tidy:           ; go mod tidy
    ```

## Error → HTTP status mapping (handler-level)

Centralized per-handler `errors.Is` switch, so `response.go` stays a dumb envelope writer:
- `ErrTodoNotFound` → 404, code `TODO_NOT_FOUND`
- `ErrInvalidInput` / validator binding failure → 400, code `INVALID_INPUT`
- anything else → 500, code `INTERNAL_ERROR`

This is the pedagogical point: usecase/repo layers never know about HTTP; only the handler translates domain errors to transport semantics.

## Suggested reading order for the junior dev

1. `domain/entity/todo.go` — what a Todo *is*, no framework noise
2. `domain/repository/todo_repository.go` — the contract, no SQL/HTTP mentioned
3. `domain/errors/errors.go` — vocabulary of failure, independent of transport
4. `usecase/interfaces/create_todo.go` (then the other 5) — just the contract for one business operation, no logic yet
5. `usecase/todo/create_todo.go` (then the other 5) — the implementation of that contract; note the constructor returns the interface type from step 4, same pattern as the repository port vs its GORM implementation
6. `infrastructure/persistence/todo_repo.go` — where SQL/GORM finally appears
7. `infrastructure/database/postgres.go` — connection setup, fail-fast
8. `delivery/http/response/response.go` — the response shape, decided once
9. `delivery/http/handlers/todo_handler.go` — where HTTP meets usecase
10. `delivery/http/middleware/*.go` — cross-cutting concerns per request
11. `delivery/http/routes.go` — URL → handler mapping
12. `infrastructure/container/container.go` — the wiring diagram (only file allowed to break layering, by design)
13. `infrastructure/server/server.go` — how the process starts listening and, on SIGINT/SIGTERM, drains in-flight requests before exiting
14. `cmd/api/main.go` — the full boot sequence in ~7 lines
15. `cmd/migrate/main.go` and `cmd/bootstrap/main.go` — the two small standalone binaries that share `config` and `container` with `cmd/api` but never run inside the HTTP server itself

Then trace one request end-to-end: `POST /api/v1/todos` → routes → handler.Create → usecase.Execute → persistence.Create → Postgres → back up to `response.Created`.

## Verification

```bash
cd /home/ebinb/to-do
cp .env.example .env        # set DB credentials for local Postgres
createdb todo_db            # or psql -c "CREATE DATABASE todo_db;"
make migrate-up             # runs cmd/migrate, no external migrate binary needed
make bootstrap               # optional: seeds a few sample todos, safe to re-run
make run                    # starts on configured port (default 8080)
```

Manual checks via curl for all 6 endpoints (create, list, get by id, update, toggle, delete) plus one 404 check (`GET /api/v1/todos/00000000-0000-0000-0000-000000000000` → expect `{"success":false,"error":{"code":"TODO_NOT_FOUND",...}}` at HTTP 404).

Unit tests: `internal/usecase/todo/create_todo_test.go` (and siblings) using `testify/mock` against a mocked `TodoRepository` — demonstrates the usecase layer is testable with zero DB/HTTP/Gin. Because handlers depend on usecase *interfaces* from `usecase/interfaces`, handler-level tests can similarly mock `interfaces.CreateTodoUsecase` etc. directly, without a real repository or database. Run via `make test`.

## Critical files
- `internal/domain/repository/todo_repository.go`
- `internal/usecase/interfaces/create_todo.go` (representative of the ports package)
- `internal/usecase/todo/create_todo.go` (representative of the implementation package)
- `internal/infrastructure/persistence/todo_repo.go`
- `internal/delivery/http/handlers/todo_handler.go`
- `internal/infrastructure/container/container.go`
- `internal/infrastructure/server/server.go`
- `internal/delivery/http/response/response.go`
- `migrations/000001_create_todos.up.sql`
- `cmd/migrate/main.go`
- `cmd/bootstrap/main.go`
- `pkg/config/config.go`
