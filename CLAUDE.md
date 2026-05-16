# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Changes (spec-driven work)

> This approach is based on [Kiro's spec methodology](https://kiro.dev/docs/specs/). The [Requirements-First workflow](https://kiro.dev/docs/specs/feature-specs/requirements-first/) is the standard used here: specify system behaviour before making architectural decisions.

Non-trivial features and bug fixes are tracked as a **change** — a folder at `.claude/changes/<change-name>/` containing up to three spec files. Use specs for anything complex, costly to get wrong, or requiring iterative design. Skip specs for exploratory/prototype work.

### Spec files

**`requirements.md`** — the *what*. Organise by feature area (H2) and user story group (H3). Each requirement uses EARS notation:

```
WHEN <condition> THE SYSTEM SHALL <action>
```

Example:
```
## Device Management

### Add device
WHEN a user submits a valid new-device form THE SYSTEM SHALL create the device record and record a creation history entry.
WHEN a user submits a device name that already exists THE SYSTEM SHALL display an "Name already taken" error.
```

Also cover edge cases and error-handling scenarios.

**`design.md`** — the *how*. Sections: system architecture and components, sequence diagrams, data models and interfaces, error-handling approach, testing strategy.

**`tasks.md`** — the *steps*. Discrete, trackable implementation tasks, each with a clear description, expected outcome, and any dependencies. Mark tasks required vs optional. Work through independent tasks first, then dependent ones in order.

**`bugfix.md`** — replaces `requirements.md` for bug fixes. Three sections using their own notation:

```
## Current Behavior (Defect)
WHEN <condition> THEN the system <incorrect behavior>

## Expected Behavior (Correct)
WHEN <condition> THEN the system SHALL <correct behavior>

## Unchanged Behavior (Regression Prevention)
WHEN <condition> THEN the system SHALL CONTINUE TO <existing behavior>
```

The "Unchanged Behavior" section is the key addition — explicitly locking down what must not change prevents regressions. The `design.md` for a bugfix includes root cause analysis; `tasks.md` includes tests that verify the bug is fixed and unchanged behavior is preserved.

### Workflow

**Feature:**
1. Create `requirements.md` and agree on it before writing `design.md`.
2. Create `design.md` and agree on it before writing `tasks.md`.
3. Execute `tasks.md` one task at a time, marking each done as you go.

**Bugfix:**
1. Create `bugfix.md` (current / expected / unchanged behavior) and agree on it.
2. Create `design.md` including root cause analysis.
3. Create and execute `tasks.md`, including tests for fix and regression prevention.

Before starting any non-trivial feature, refactor, or bug fix, check `.claude/changes/` for an existing change folder. If none exists, create one and start with `requirements.md` (feature) or `bugfix.md` (bug).

## What is KeyBook

KeyBook is a web app for managing devices, persons, and properties, with automatic audit history for all changes. The frontend is a SvelteKit static site; the backend is a Go binary that embeds PocketBase (SQLite + REST API). The compiled frontend is served by the backend from `backend/internal/frontend/build/`.

## Architecture

```
frontend/          SvelteKit (Svelte 5, TypeScript, Carbon Design System)
backend/
  cmd/keybook.go   Entry point — wires DI container, registers PocketBase hooks
  internal/
    repositories/  Data access layer (PocketBase DAO queries)
    services/      Business logic; history-tracking services called by hooks
    dtos/          Data transfer objects for all entities
    helpers/       PocketBase DAO error utilities
    frontend/build/ Gitignored — populated from frontend build output
database/
  pb_schema.json   PocketBase collection definitions
```

### Key patterns

**Frontend:** Business logic lives in `.svelte.ts` module files under `src/lib/modules/` (one per entity: device, person, property, persondevice, user). These implement typed interfaces and use Svelte 5 reactive primitives (`$state`, `$derived.by`). Shared state is distributed via Svelte's context API (set in the user layout, consumed via `getContext()`), not stores.

**Backend:** Dependency injection via `go.uber.org/dig`. Repository pattern for data access; service layer for business logic. Audit history is recorded automatically — PocketBase `OnModelAfterCreate` and `OnModelBeforeUpdate` hooks call history services for every entity type.

**Frontend → backend:** The PocketBase JS SDK (`pocketbase` npm package) is the only HTTP client. Base URL comes from the `PUBLIC_POCKETBASE_URL` env var. Auth state is persisted in cookies via `src/lib/api/backend-client.ts`. The user layout (`src/routes/user/+layout.ts`) guards all `/user/*` routes and redirects to `/auth` if unauthenticated.

## Layered Architecture

All code must follow a layered architecture with clear separation of concerns. Each layer may only depend on the layer directly below it.

**Backend:**

| Layer | KeyBook implementation |
|---|---|
| **Hooks** | PocketBase event handlers in `cmd/keybook.go` — route events to services |
| **Services** | Business logic and audit history recording (`internal/services/`) |
| **Repositories** | PocketBase DAO queries (`internal/repositories/`) |

**Frontend:**

| Layer | KeyBook implementation |
|---|---|
| **Components** | `.svelte` files — UI and user interaction |
| **Modules** | `.svelte.ts` files in `src/lib/modules/` — application state and business logic |
| **SDK** | PocketBase JS SDK — all HTTP and data access |

Components call module methods; modules call the SDK directly. No component accesses the SDK directly.

### References

1. [Microsoft Learn — Infrastructure Persistence Layer Design](https://learn.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/infrastructure-persistence-layer-design)
2. [From Request to Database: Three-Layer Architecture in API Development](https://konstantinmb.medium.com/from-request-to-database-understanding-the-three-layer-architecture-in-api-development-1c44c973c7af)
3. [Repository Pattern with Layered Architecture](https://medium.com/@leadcoder/repository-pattern-with-layered-architecture-35f7b9182ebf)
4. [Layered Architecture Template for REST APIs (Java/Spring Boot)](https://kamilmazurek.pl/layered-architecture-template)
5. [Service–Repository Pattern in Action](https://medium.com/@albinaji.official/service-repository-pattern-in-action-0db4bb9a474b)
6. [Repository and Services Pattern in a Multilayered Architecture](https://www.vodovnik.com/repository-and-services-pattern-in-a-multilayered-architecture/)
7. [Clean Architecture — Incorporating Repository Pattern](https://medium.com/@bert.oneill/clean-architecture-incorporating-repository-pattern-388742e0b54e)
8. [Unpacking Clean Architecture Layers — Domain, Application, Infrastructure Services](https://www.dandoescode.com/blog/unpacking-the-layers-of-clean-architecture-domain-application-and-infrastructure-services)
9. [Clean Architecture: Understanding the Infrastructure and Persistence Layers](https://dev.to/moh_moh701/what-is-clean-architecture-understanding-the-infrastructure-and-persistence-layers-2pca)
10. [clean-architecture-api-boilerplate (GitHub)](https://github.com/luizomf/clean-architecture-api-boilerplate/blob/master/README.md)
11. [Building a Layered Architecture in NestJS & TypeScript](https://medium.com/@patrick.cunha336/building-a-layered-architecture-in-nestjs-typescript-repository-pattern-dtos-and-validators-08907a8ac4cb) — **⚠️ Uses NestJS/Node.js. This project uses SvelteKit for the frontend — apply the structural concepts only, not the framework specifics.**
12. [Clean Architecture Design Guide for Backend API Developments](https://naskay.com/blog/clean-architecture-design-guide-for-backend-api-developments/)

## Development rules

- Always read a file before editing it.
- Schema changes go through `database/pb_schema.json` only — never edit the PocketBase database directly.
- No speculative abstractions — only build what is needed now.
- Format Go code with `gofmt` and `goimports` before committing.
- Run `npm run lint` in `frontend/` to check TypeScript/Svelte style before committing.
- Code style: [Google Go style guide](https://google.github.io/styleguide/go/guide), [Google TypeScript style guide](https://google.github.io/styleguide/tsguide.html), [Svelte style guide](https://svelte.dev/docs/svelte/style-guide).

## Testing

### Mandate

**All tests must be run on the OpenBSD server** (see `CLAUDE.local.md` for connection details). Do not run tests locally.

**Unit tests must be written before the implementation code they cover (TDD).** Write the test, watch it fail, then write the minimum code to make it pass.

### Test types

| Type | Scope | When required |
|---|---|---|
| **Unit** | Single function, module, or component in isolation | Always — written first |
| **Integration** | Multiple components or layers working together (e.g. repository + service, component + store) | Always for non-trivial interactions |
| **E2E** | Full user flow through the running app via browser | Always for user-facing features |
| **Contract** | API shape between frontend and backend (request/response structure) | When adding or changing PocketBase collection endpoints |

### Frontend tests

Every integration test exercises Svelte components together with their modules against a real running PocketBase instance — no SDK mocking.

| Command | Tool | What it covers |
|---|---|---|
| `npm run test:unit` | Vitest | Unit and integration tests in `src/` |
| `npm run test:integration` | Playwright | E2E flows against the running app |

Place unit/integration test files alongside the code they test (`*.test.ts` or `*.spec.ts`). E2E tests live in `tests/`.

### Backend tests

Every integration test starts a real PocketBase HTTP server using a `t.TempDir()` data directory — no database mocking.

```sh
go test ./...               # all tests
go test ./internal/...      # specific package tree
```

Place test files alongside the Go source (`*_test.go`). Use table-driven tests for repository and service logic.

### In specs

`tasks.md` must include explicit test tasks. For features: unit tests as the first task for each component, followed by integration and E2E tasks. For bugfixes: tasks must include a test that reproduces the bug before it is fixed, and regression tests drawn from `bugfix.md`'s "Unchanged Behavior" section.

## Frontend commands

Run from `frontend/`:

```sh
npm run dev          # dev server — address in CLAUDE.local.md
npm run build        # production build → frontend/build/
npm run check        # svelte-check + TypeScript
npm run lint         # prettier + eslint (check only)
npm run format       # prettier --write
npm run test:unit    # vitest
npm run test:integration  # playwright
```

Code style: tabs, single quotes, 100-char line width (Prettier). TypeScript strict mode.

## Backend commands

Run from `backend/`:

```sh
go build -o build/keybook ./cmd/     # compile
./build/keybook serve --dev --http <address>  # address in CLAUDE.local.md
go test ./...                         # tests
```

## Deploying to the OpenBSD server

Server addresses, credentials, and deployment steps are in `CLAUDE.local.md` (gitignored).

**Key facts (non-sensitive):**
- `sshpass` is installed in Cygwin — always use `/c/cygwin64/bin/sshpass` alongside Cygwin's `ssh`/`sftp`; it is not on the Git Bash PATH.
- Use `git -c http.sslVerify=false push` on Windows if SSL verification errors occur.

@CLAUDE.local.md
