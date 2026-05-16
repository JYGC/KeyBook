# Design — full-collections

## 1. Overview

This change restructures the data model from a flat, property-centric design to a normalised multi-role design. Persons become first-class entities (not property-scoped). Ownership, tenancy, and management are expressed through dedicated relation collections rather than embedded fields. A new item registry supports physical household items and entry devices. Firm (cobrand) ownership and agent management are introduced as first-class concepts.

The old database is dropped and replaced from scratch. Columns from old collections are carried forward onto their counterpart collections in the new schema where they are not already covered. The new schema is defined in `database/pb_schema.json` and imported via PocketBase's schema import on a fresh database — consistent with the development rule in `CLAUDE.md`.

All new backend and frontend code follows the layered architecture defined in `CLAUDE.md`: API → Application → Service → Repository → Store. This change introduces `internal/application/` (backend) and `src/lib/services/`, `src/lib/repositories/` (frontend) as new layers.

## 2. Schema changes

### 2.1 Complete collection definitions

All columns listed. Columns marked *(carried over)* come from the old schema. History collections follow the existing pattern (snapshot JSON + `statedDateTime` per entity) and are retained for `persons` and `properties` only — relation-only tables rely on PocketBase's built-in `created`/`updated` timestamps.

---

**`users`** *(auth — created by PocketBase by default; only field additions shown)*
| Field | Type | Notes |
|---|---|---|
| `name` | text | |
| `avatar` | file | *(carried over)* |

---

**`persons`**
| Field | Type | Notes |
|---|---|---|
| `name` | text, required | |
| `DOB` | date | |
| `user` | relation → users | optional |
| `profileImage` | file | *(carried over)* |

---

**`properties`**
| Field | Type | Notes |
|---|---|---|
| `address` | text, required | |

---

**`propertyOwners`**
| Field | Type | Notes |
|---|---|---|
| `property` | relation → properties, required | |

---

**`personPropertyOwners`**
| Field | Type | Notes |
|---|---|---|
| `person` | relation → persons, required | |
| `propertyOwner` | relation → propertyOwners, required | |
| `name` | text | display name for this ownership |

---

**`cobrands`**
| Field | Type | Notes |
|---|---|---|
| `name` | text, required | firm/company name |

---

**`cobrandAdmins`**
| Field | Type | Notes |
|---|---|---|
| `userId` | relation → users, required | |
| `cobrand` | relation → cobrands, required | |

---

**`cobrandPropertyManagers`**
| Field | Type | Notes |
|---|---|---|
| `cobrand` | relation → cobrands, required | |
| `property` | relation → properties, required | |

---

**`cobrandPropertyOwners`**
| Field | Type | Notes |
|---|---|---|
| `cobrand` | relation → cobrands, required | |
| `propertyOwner` | relation → propertyOwners, required | |
| `name` | text | display name for this ownership |

---

**`agents`**
| Field | Type | Notes |
|---|---|---|
| `person` | relation → persons, required | |
| `cobrand` | relation → cobrands, required | |

---

**`propertyAgents`**
| Field | Type | Notes |
|---|---|---|
| `agent` | relation → agents, required | |
| `property` | relation → properties, required | |

---

**`households`**
| Field | Type | Notes |
|---|---|---|
| `person` | relation → persons, required | |
| `property` | relation → properties, required | |

---

**`tenants`**
| Field | Type | Notes |
|---|---|---|
| `person` | relation → persons, required | |
| `property` | relation → properties, required | |

---

**`items`**
| Field | Type | Notes |
|---|---|---|
| `name` | text, required | |
| `description` | text | |
| `picture` | file | (was `image` on old `devices`) |

---

**`entryDevices`**
| Field | Type | Notes |
|---|---|---|
| `item` | relation → items, required | |
| `deviceType` | select (Fob/Key/Remote/RoomKey/MailboxKey), required | (was `type` on old `devices`) |
| `identifier` | text, required | *(carried over from `devices`)* |
| `defunctReason` | select (None/Lost/Damaged/Retired/Stolen), required | *(carried over from `devices`)* |

---

**`propertyItems`**
| Field | Type | Notes |
|---|---|---|
| `item` | relation → items, required | |
| `property` | relation → properties, required | |

---

**`personItems`**
| Field | Type | Notes |
|---|---|---|
| `person` | relation → persons, required | |
| `item` | relation → items, required | |

---

**`personHistory`** *(unchanged)*
| Field | Type |
|---|---|
| `personId` | text, required |
| `person` | json (snapshot), required |
| `statedDateTime` | date |
| `description` | text |
| `property` | relation → properties, required |

---

**`propertyHistory`** *(unchanged)*
| Field | Type |
|---|---|
| `propertyId` | text, required |
| `property` | json (snapshot), required |
| `statedDateTime` | date |
| `description` | text |

---

### 2.2 Collections not carried forward

The following old collections are dropped with the database and have no equivalent in the new schema:

- `devices` — superseded by `items` + `entryDevices`
- `personDevices` — superseded by `personItems`
- `deviceHistory` — dropped (entry devices tracked via `entryDevices` timestamps)
- `personDeviceHistory` — dropped

### 2.3 Access rule implications

Current rules reference `property.owners.id`. After this change, ownership is resolved through:
```
propertyOwners.property = <propertyId>
→ personPropertyOwners.propertyOwner
→ person.user.id = @request.auth.id
```
OR
```
cobrandPropertyOwners.propertyOwner
→ cobrand
→ cobrandAdmins.userId = @request.auth.id
```
Access rules on all collections will need updating to use the new ownership chain.

## 3. Collection relationship diagram

```
users ──────────────────────────────────────── cobrandAdmins ── cobrands
  │                                                                  │
  └── persons                                          cobrandPropertyManagers ── properties
        │                                              cobrandPropertyOwners ─────── │
        ├── personPropertyOwners ── propertyOwners ─────────────────────────── properties
        ├── agents ──────────────── propertyAgents ── properties
        ├── households ──────────────────────────── properties
        ├── tenants ─────────────────────────────── properties
        └── personItems ── items ── propertyItems ── properties
                                 └── entryDevices
```

## 4. Backend changes

### 4.1 New repositories (`internal/repositories/`)

One repository per new collection following the existing pattern:

- `CobrandRepository`
- `CobrandAdminRepository`
- `CobrandPropertyManagerRepository`
- `CobrandPropertyOwnerRepository`
- `PropertyOwnerRepository`
- `PersonPropertyOwnerRepository`
- `AgentRepository`
- `PropertyAgentRepository`
- `HouseholdRepository`
- `TenantRepository`
- `ItemRepository`
- `PropertyItemRepository`
- `PersonItemRepository`
- `EntryDeviceRepository`

### 4.2 Modified repositories

- `PersonRepository` — update queries to remove `type`/`property` fields; add `user`, `DOB`.
- `PropertyRepository` — remove `owners`/`managers` field references.
- Delete `DeviceRepository`, `PersonDeviceRepository`.

### 4.3 New application services (`internal/application/`)

One application service per entity group, coordinating repository calls without containing business logic. Each is injected with its required repositories via the `dig` container.

- `ItemApplicationService` — orchestrates item CRUD, entry device designation, and person/property item associations
- `CobrandApplicationService` — orchestrates cobrand CRUD, admin assignment, and property manager assignment
- `AgentApplicationService` — orchestrates agent registration and property assignment
- `PropertyOwnerApplicationService` — orchestrates property owner creation and person/cobrand owner linkage
- `PersonApplicationService` — orchestrates person CRUD and user account linking
- `PropertyApplicationService` — orchestrates property CRUD

Hook handlers in `cmd/keybook.go` call application services; application services call service-layer functions; services call repositories.

### 4.4 Services and hooks

No new history services are required for the relation collections. The existing `PersonHistoryServices` and `PropertyHistoryServices` hooks continue to apply. Remove `DeviceHistoryServices` and `PersonDeviceHistoryServices` after migration.

### 4.5 DTOs (`internal/dtos/`)

Add DTOs for each new collection. Update `PersonDtos` and `PropertyDtos` to reflect field changes.

## 5. Frontend changes

### 5.1 New routes

| Route | Purpose |
|---|---|
| `/user/items/` | Item list, add, edit |
| `/user/items/[id]/` | Item detail with entry device status |
| `/user/cobrands/` | Cobrand list, add |
| `/user/cobrands/[id]/` | Cobrand detail (admins, agents, managed/owned properties) |
| `/user/agents/` | Agent list, add |
| `/user/agents/[id]/` | Agent detail with assigned properties |

### 5.2 New application modules (`src/lib/modules/`)

- `item/` — `ItemListModule`, `ItemDetailModule`
- `cobrand/` — `CobrandListModule`, `CobrandDetailModule`
- `agent/` — `AgentListModule`, `AgentDetailModule`

Modules are the Application layer: they orchestrate use cases and hold reactive state via Svelte 5 primitives. They call services; they do not call repositories or the SDK directly.

### 5.3 Updated modules

- `person/` — remove type/property handling; add user-link, DOB, role derivation via relation collections.
- `property/` — remove direct owners/managers; derive ownership from new collections.

### 5.4 New contexts

- `ItemContext` (following existing pattern)
- `CobrandContext`
- `AgentContext`

### 5.5 New services (`src/lib/services/`)

One service per entity group, containing business logic and validation:

- `ItemService` — item validation, entry device state transitions, defunct reason validation
- `CobrandService` — cobrand validation, admin uniqueness enforcement
- `AgentService` — agent registration rules, property assignment validation
- `PropertyOwnerService` — ownership validation, duplicate owner prevention
- `PersonService` — person validation, role derivation from relation collections
- `PropertyService` — property validation, ownership chain traversal

Services call repositories; they do not call the SDK directly.

### 5.6 New repositories (`src/lib/repositories/`)

One repository per collection, abstracting all PocketBase SDK calls:

- `ItemRepository`, `EntryDeviceRepository`, `PropertyItemRepository`, `PersonItemRepository`
- `CobrandRepository`, `CobrandAdminRepository`, `CobrandPropertyManagerRepository`, `CobrandPropertyOwnerRepository`
- `AgentRepository`, `PropertyAgentRepository`
- `PropertyOwnerRepository`, `PersonPropertyOwnerRepository`
- `PersonRepository`, `PropertyRepository`

Repositories are the only layer that calls the PocketBase JS SDK.

## 6. Error-handling approach

- **Backend repositories**: Translate PocketBase DAO errors to domain errors using helpers in `internal/helpers/`. Surface HTTP status and message on failure.
- **Backend application services**: Propagate repository errors without wrapping unless additional context is required.
- **Backend hook handlers**: Log errors and return appropriate PocketBase error responses to the caller.
- **Frontend repositories**: Propagate PocketBase SDK errors (with HTTP status codes) to the service layer unchanged.
- **Frontend services**: Translate repository errors into user-facing messages where appropriate.
- **Frontend modules**: Expose error state as reactive `$state` fields; components display errors without containing error-handling logic.

## 7. Testing strategy

- **Unit tests** (written first per TDD): Cover DTOs (field mapping, required fields) and service business logic in isolation.
- **Integration tests (backend)**: Each test starts a real PocketBase HTTP server via `t.TempDir()` — no database mocking. Cover repository CRUD, application service orchestration, and hook event flows.
- **Integration tests (frontend)**: Exercise Svelte components together with their modules, services, and repositories against a real running PocketBase instance — no SDK mocking.
- **E2E tests**: Playwright tests cover full user flows for each new route.
- **Contract tests**: Verify PocketBase collection access rules (HTTP status codes, auth enforcement) via direct HTTP requests. Written first, before schema is imported.
- All tests run on the OpenBSD server (see `CLAUDE.local.md`).

## 8. Data migration path

1. Update `database/pb_schema.json` with all new collections as specified in section 2.1.
2. On the OpenBSD server: stop the binary, delete `pb_data/`, import the new schema via PocketBase's admin UI or `--importcollections` flag, restart.
3. Migrate existing data (run as a one-shot script):
   - For each `devices` record → create `items` + `entryDevices` record.
   - For each `personDevices` record → create `personItems` record.
   - For each `persons` record with type `Owner` → create `propertyOwners` + `personPropertyOwners` record.
   - For each `persons` record with type `Tenant` → create `tenants` record.
   - For each `persons` record with type `Agent` → create `agents` + `propertyAgents` record.
   - For each `persons` record with type `Household` → create `households` record.
   - For each `properties` record → create `propertyOwners` records from the `owners` relation.
4. Update all PocketBase access rules to use the new ownership chain.
