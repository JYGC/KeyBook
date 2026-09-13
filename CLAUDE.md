# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

It covers **how** work is done here: the spec-driven approach to changes, the architecture standard all code follows, and the coding and testing rules. Everything specific to this repository — what the app is, its layout, its commands, and where things run — lives in @CLAUDE-project.md.

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
## Account Management

### Register account
WHEN a user submits a valid registration form THE SYSTEM SHALL create the account and record a creation history entry.
WHEN a user submits an email address that is already registered THE SYSTEM SHALL display an "Email already taken" error.
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

## Layered Architecture

All code must follow a layered architecture with clear separation of concerns. Each layer may only depend on the layer directly below it.

| Layer | Responsibility |
|---|---|
| **API** | HTTP handlers, request/response mapping, input validation |
| **Application** | Use-case orchestration; coordinates services without containing business logic |
| **Service** | Business logic and domain rules |
| **Repository** | Data access abstraction; hides persistence details from services |
| **Store** | Persistence (database queries, external API calls) |

Optional: **Domain** (pure entities and value objects, no dependencies), **DTO/Schema** (typed data transfer objects at layer boundaries).

See @CLAUDE-project.md for how these layers map onto this repository's frontend and backend directories.

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
- Never edit a database directly — schema changes go through migrations.
- No speculative abstractions — only build what is needed now.
- Prefer long, self-documenting names over comments. Methods, functions, classes, variables, types, and every other namable thing should carry their documentation in the name — `recordDeviceOwnershipTransfer` rather than `transfer` with a comment explaining what is transferred. Reach for a comment only when the *why* cannot be expressed in a name (non-obvious constraints, external quirks, trade-offs).
- Format and lint code before committing, using the project's configured tooling.
- Code style: [Google Go style guide](https://google.github.io/styleguide/go/guide), [Google TypeScript style guide](https://google.github.io/styleguide/tsguide.html), [Svelte style guide](https://svelte.dev/docs/svelte/style-guide).
- Frontend formatting is Prettier-enforced: tabs, single quotes, 100-char line width. TypeScript runs in strict mode.

## Testing

### Mandate

**Unit tests must be written before the implementation code they cover (TDD).** Write the test, watch it fail, then write the minimum code to make it pass.

**Integration tests run against real dependencies** — a real server, a real database instance. No mocking of the database, the SDK, or the API client.

Test commands and where they must be run are in @CLAUDE-project.md.

### Test types

| Type | Scope | When required |
|---|---|---|
| **Unit** | Single function, module, or component in isolation | Always — written first |
| **Integration** | Multiple components or layers working together (e.g. repository + service, component + store) | Always for non-trivial interactions |
| **E2E** | Full user flow through the running app via browser | Always for user-facing features |
| **Contract** | API shape between frontend and backend (request/response structure) | When adding or changing collection endpoints |

### Conventions

- Place unit and integration test files alongside the code they test.
- Use table-driven tests for repository and service logic.

### In specs

`tasks.md` must include explicit test tasks. For features: unit tests as the first task for each component, followed by integration and E2E tasks. For bugfixes: tasks must include a test that reproduces the bug before it is fixed, and regression tests drawn from `bugfix.md`'s "Unchanged Behavior" section.
