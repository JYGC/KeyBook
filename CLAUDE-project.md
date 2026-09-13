# CLAUDE-project.md

Repository-specific context for KeyBook: what it is, how it is laid out, and how to build, run, and test it. General engineering standards live in [CLAUDE.md](./CLAUDE.md).

## What is KeyBook

KeyBook is a web app for managing devices, persons, and properties, with automatic audit history for all changes. The frontend is a SvelteKit static site served independently; the backend is a Go binary that embeds PocketBase (SQLite + REST API).

## Repository layout

```
frontend/
  src/
    lib/
      modules/      Application layer — use-case orchestration, reactive state
      services/     Service layer — business logic (to be introduced)
      repositories/ Repository layer — PocketBase SDK abstraction (to be introduced)
    routes/         SvelteKit pages and layouts
backend/
  cmd/keybook.go    API layer — entry point, DI wiring, PocketBase hook handlers
  migrations/       Schema source of truth — PocketBase migrations
  internal/
    application/    Application layer — use-case orchestration (to be introduced)
    services/       Service layer — business logic and audit history
    repositories/   Repository layer — PocketBase DAO queries
    dtos/           DTO layer — data transfer objects at layer boundaries
    helpers/        PocketBase DAO error utilities
```

## Layer mapping

The codebase is being migrated toward the layered architecture defined in [CLAUDE.md](./CLAUDE.md#layered-architecture). New code must follow the target patterns; existing code is updated incrementally.

**Backend:**

| Layer | KeyBook target |
|---|---|
| **API** | PocketBase hook handlers in `cmd/keybook.go` |
| **Application** | `internal/application/` — to be introduced |
| **Service** | `internal/services/` |
| **Repository** | `internal/repositories/` |
| **Store** | PocketBase DAO |
| **DTO** | `internal/dtos/` |

**Frontend:**

| Layer | KeyBook target |
|---|---|
| **API** | N/A — frontend consumes the PocketBase REST API via the SDK |
| **Application** | `.svelte.ts` modules in `src/lib/modules/` |
| **Service** | `src/lib/services/` — to be introduced |
| **Repository** | `src/lib/repositories/` — to be introduced |
| **Store** | PocketBase JS SDK |

## Key patterns

**Frontend (target):** Components (`.svelte`) handle UI and user interaction only — no business logic. Application modules (`.svelte.ts` in `src/lib/modules/`) orchestrate use cases and hold reactive state via Svelte 5 primitives (`$state`, `$derived.by`). Service layer (`src/lib/services/`) contains business logic. Repository layer (`src/lib/repositories/`) abstracts all PocketBase SDK calls. State is distributed via Svelte's context API (set in the user layout, consumed via `getContext()`), not stores.

**Backend (target):** Dependency injection via `go.uber.org/dig`. Hook handlers (`cmd/keybook.go`) are the API layer — they receive PocketBase events and delegate to the application layer. Application layer (`internal/application/`) orchestrates services per use case without containing business logic. Service layer (`internal/services/`) contains business logic and audit history recording. Repository layer (`internal/repositories/`) abstracts all PocketBase DAO access. DTOs (`internal/dtos/`) cross layer boundaries.

**Frontend → backend:** The PocketBase JS SDK (`pocketbase` npm package) is the only HTTP client (Store layer). Base URL comes from the `PUBLIC_POCKETBASE_URL` env var. Auth state is persisted in cookies via `src/lib/api/backend-client.ts`. The user layout (`src/routes/user/+layout.ts`) guards all `/user/*` routes and redirects to `/auth` if unauthenticated.

## Project rules

- Schema changes go through PocketBase migrations in `backend/migrations/` only — never edit the database directly.
- Run `npm run lint` in `frontend/` to check TypeScript/Svelte style before committing.
- Format Go code with `gofmt` and `goimports` before committing.

## Testing

All tests must be run on the OpenBSD server — use `/openbsd-run --test backend` and `/openbsd-run --test frontend`, which carry the connection details and the exact invocations. Do not run tests locally.

### Frontend tests

Every integration test exercises Svelte components together with their modules against a real running PocketBase instance — no SDK mocking.

| Command | Tool | What it covers |
|---|---|---|
| `npm run test:unit` | Vitest | Unit and integration tests in `src/` |
| `npm run test:integration` | Playwright | E2E flows against the running app |

Place unit/integration test files alongside the code they test (`*.test.ts` or `*.spec.ts`). E2E tests live in `tests/`.

Both suites talk to a **live** dev backend — start it first. Vitest defaults to watch mode, so pass `-- --run` when invoking it non-interactively. The Playwright suite pins `channel: 'msedge'` and cannot run on the OpenBSD server; run it from a workstation on the same subnet.

### Backend tests

Every integration test starts a real PocketBase HTTP server using a `t.TempDir()` data directory — no database mocking.

```sh
go test ./...               # all tests
go test ./internal/...      # specific package tree
```

Place test files alongside the Go source (`*_test.go`).

## Frontend commands

Run from `frontend/`:

```sh
npm run dev          # dev server — binds 192.168.8.144:5173 per vite.config.ts
npm run build        # production build → frontend/build/
npm run check        # svelte-check + TypeScript
npm run lint         # prettier + eslint (check only)
npm run format       # prettier --write
npm run test:unit    # vitest
npm run test:integration  # playwright
```

## Backend commands

Run from `backend/`:

```sh
go build -o build/keybook ./cmd/                     # compile
./build/keybook serve --dev --http 192.168.8.144:8090  # dev server
go test ./...                                        # tests
```

PocketBase is v0.22, so the admin CLI verb is `admin` (`./build/keybook admin create <email> <password>`), not `superuser`. Migrations registered in `backend/migrations/` run automatically on `serve`.

## Running on the OpenBSD server

Everything — start, stop, status, build, test, lint, check, dependency installs — goes through the `/openbsd-run` command (`.claude/commands/openbsd-run.md`, untracked). It holds the connection details, the verified command forms, and the traps worth knowing about. Read it before driving the server by hand.

**Key facts:**
- The dev stack runs as `junying` from `/home/junying/Source/KeyBook` and binds `192.168.8.144` — backend on `:8090`, Vite dev server on `:5173`. SSH is key auth to `192.168.8.145`; no password helper is involved.
- `192.168.8.143` is the **same machine** but hosts another project's production PocketBase, also on port `:8090`. Only the address separates them, so never bind a wildcard address and never signal a process you do not own.
- The frontend's effective backend URL comes from `frontend/.env.local` on the server (untracked); the tracked `frontend/.env` points at a dead port. An inline `PUBLIC_POCKETBASE_URL=…` overrides both, for `npm run dev` and `npm run build` alike.
