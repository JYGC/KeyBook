# Design — full-collections

## 1. Overview

This change restructures the data model from a flat, property-centric design to a normalised multi-role design. Persons become first-class entities (not property-scoped). Ownership, tenancy, and management are expressed through dedicated relation collections rather than embedded fields. A new item registry supports physical household items and entry devices. Firm (cobrand) ownership and agent management are introduced as first-class concepts.

The old database is dropped and replaced from scratch. Columns from old collections are carried forward onto their counterpart collections in the new schema where they are not already covered. The new schema is defined and managed exclusively via PocketBase migrations in `backend/migrations/`. `database/pb_schema.json` is deleted.

All new backend and frontend code follows the layered architecture defined in `CLAUDE.md`: API → Application → Service → Repository → Store. This change introduces `internal/application/` (backend) and `src/lib/services/`, `src/lib/repositories/` (frontend) as new layers.

## 2. Schema changes

### 2.1 Complete collection definitions

All columns listed. Columns marked *(carried over)* come from the old schema. Three history collections (`propertyHistories`, `tenantHistories`, `propertyItemHistories`) track property-scoped changes — relation-only tables rely on PocketBase's built-in `created`/`updated` timestamps.

---

**`users`** *(auth — created by PocketBase by default; `email` is kept; `name` and `avatar` are included by default but must be removed)*

---

**`persons`**
| Field | Type | Notes |
|---|---|---|
| `name` | text, required | |
| `DOB` | date, required | |
| `user` | relation → users, unique | optional (null allowed; non-null values must be unique) |
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

**`personPropertyOwners`** *(unique on `person` + `propertyOwner`)*
| Field | Type | Notes |
|---|---|---|
| `person` | relation → persons, required | |
| `propertyOwner` | relation → propertyOwners, required | |

---

**`cobrands`**
| Field | Type | Notes |
|---|---|---|
| `name` | text, required | firm/company name |

---

**`cobrandAdmins`** *(unique on `user` + `cobrand`)*
| Field | Type | Notes |
|---|---|---|
| `user` | relation → users, required | |
| `cobrand` | relation → cobrands, required | |

---

**`cobrandPropertyManagers`** *(unique on `cobrand` + `property`)*
| Field | Type | Notes |
|---|---|---|
| `cobrand` | relation → cobrands, required | |
| `property` | relation → properties, required | |

---

**`cobrandPropertyOwners`** *(unique on `cobrand` + `propertyOwner`)*
| Field | Type | Notes |
|---|---|---|
| `cobrand` | relation → cobrands, required | |
| `propertyOwner` | relation → propertyOwners, required | |

---

**`agents`** *(unique on `person` + `cobrand`)*
| Field | Type | Notes |
|---|---|---|
| `person` | relation → persons, required | |
| `cobrand` | relation → cobrands, required | |

---

**`propertyAgents`** *(unique on `agent` + `property`)*
| Field | Type | Notes |
|---|---|---|
| `agent` | relation → agents, required | |
| `property` | relation → properties, required | |

---

**`households`** *(unique on `person` + `property`)*
| Field | Type | Notes |
|---|---|---|
| `person` | relation → persons, required | |
| `property` | relation → properties, required | |

---

**`tenants`** *(unique on `person` + `property`)*
| Field | Type | Notes |
|---|---|---|
| `person` | relation → persons, required | |
| `property` | relation → properties, required | |

---

**`items`**
| Field | Type | Notes |
|---|---|---|
| `name` | text, required | |
| `description` | text, required | |
| `picture` | file | (was `image` on old `devices`) |

---

**`entryDevices`** *(unique on `item`)*
| Field | Type | Notes |
|---|---|---|
| `item` | relation → items, required | |
| `deviceType` | select (Fob/Key/Remote/RoomKey/MailboxKey), required | (was `type` on old `devices`) |
| `identifier` | text, required | *(carried over from `devices`)* |
| `defunctReason` | select (None/Lost/Damaged/Retired/Stolen), required | *(carried over from `devices`)* |

---

**`propertyItems`** *(unique on `item` + `property`)*
| Field | Type | Notes |
|---|---|---|
| `item` | relation → items, required | |
| `property` | relation → properties, required | |

---

**`personItems`** *(unique on `person` + `item`)*
| Field | Type | Notes |
|---|---|---|
| `person` | relation → persons, required | |
| `item` | relation → items, required | |

---

**`propertyHistories`**

All relation fields are nullable with cascade delete disabled — history records survive property deletion. The `snapshot` preserves readable state after the property record is gone.

| Field | Type | Notes |
|---|---|---|
| `property` | relation → properties | optional; no cascade delete; becomes null if property deleted |
| `snapshot` | json, required | full property state captured at time of change |
| `description` | text | |
| `statedDateTime` | date | |

---

**`tenantHistories`**

All relation fields are nullable with cascade delete disabled. `personSnapshot` and `propertySnapshot` capture the full state of each referenced entity at write time so the record remains readable after either is deleted.

| Field | Type | Notes |
|---|---|---|
| `property` | relation → properties | optional; no cascade delete; becomes null if property deleted |
| `propertySnapshot` | json, required | full property state captured at time of event |
| `person` | relation → persons | optional; no cascade delete; becomes null if person deleted |
| `personSnapshot` | json, required | full person state captured at time of event |
| `action` | select (Added/Removed), required | |
| `statedDateTime` | date | |

---

**`propertyItemHistories`**

All relation fields are nullable with cascade delete disabled. `itemSnapshot` and `propertySnapshot` capture the full state of each referenced entity at write time so the record remains readable after either is deleted.

| Field | Type | Notes |
|---|---|---|
| `property` | relation → properties | optional; no cascade delete; becomes null if property deleted |
| `propertySnapshot` | json, required | full property state captured at time of event |
| `item` | relation → items | optional; no cascade delete; becomes null if item deleted |
| `itemSnapshot` | json, required | full item state captured at time of event |
| `action` | select (Added/Removed), required | |
| `statedDateTime` | date | |

---

### 2.2 Collections not carried forward

The following old collections are dropped with the database and have no equivalent in the new schema:

- `devices` — superseded by `items` + `entryDevices`
- `personDevices` — superseded by `personItems`
- `deviceHistory` — dropped (entry devices tracked via `entryDevices` timestamps)
- `personDeviceHistory` — dropped
- `personHistory` — dropped; person changes are not tracked in the new schema
- `propertyHistory` — dropped; replaced by `propertyHistories`, `tenantHistories`, and `propertyItemHistories`

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
→ cobrandAdmins.user = @request.auth.id
```
Access rules on all collections will need updating to use the new ownership chain.

#### `persons` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | `@request.auth.id != ""` |
| update | `user.id = @request.auth.id` |
| delete | `user.id = @request.auth.id` |

**list / view** — a person is visible to:

```
user.id = @request.auth.id
|| tenants.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| tenants.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| tenants.property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
|| personItems.item.propertyItems.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| personItems.item.propertyItems.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| personItems.item.propertyItems.property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
```

The three groups this covers:
- **Own user** — the person whose `user` field matches the authenticated account.
- **Tenant property owners/managers** — any user who is an owner (via `personPropertyOwners` or `cobrandPropertyOwners`) or a cobrand manager (via `cobrandPropertyManagers`) of a property where this person appears in `tenants`.
- **Item property owners/managers** — any user who owns or manages a property that has an item (via `propertyItems`) currently held by this person (via `personItems`).

**create** is open to any authenticated user so that property owners can register tenants and household members who do not yet have their own user account. The person becomes visible to the creating owner once linked as a tenant (or via a held item) — the frontend receives the created record ID directly from the create response and must use it immediately to create the association record without re-querying.

Update and delete are restricted to the user account linked to that person (`persons.user`). Property owners manage their relationship to a person via the association collections (`tenants`, `households`, etc.), not by editing the person record itself.

#### `properties` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | `@request.auth.id != ""` |
| update | *(see below)* |
| delete | *(see below — same as update)* |

**list / view** — a property is visible to any user with a current relationship to it:

```
propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
|| tenants.person.user.id = @request.auth.id
|| households.person.user.id = @request.auth.id
|| propertyAgents.agent.person.user.id = @request.auth.id
```

The six groups this covers:
- **Person owners** — users linked via `personPropertyOwners → propertyOwners`.
- **Cobrand owners** — cobrand admins of cobrands linked via `cobrandPropertyOwners → propertyOwners`.
- **Cobrand managers** — cobrand admins of cobrands linked via `cobrandPropertyManagers`.
- **Tenants** — users whose person record appears in `tenants` for this property.
- **Household members** — users whose person record appears in `households` for this property.
- **Agents** — users whose person record is linked via `propertyAgents → agents` for this property.

**update / delete** — restricted to owners only:

```
propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
```

Cobrand managers, tenants, household members, and agents can view but not modify the property record itself.

#### `propertyOwners` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | `@request.auth.id != ""` |
| update | *(see below)* |
| delete | *(see below — same as update)* |

**list / view** — an ownership record is visible to anyone who owns or manages the linked property:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
```

The three groups this covers:
- **Person owners** — co-owners can see each other's ownership records on the same property.
- **Cobrand owners** — cobrand admins of cobrands that co-own the property.
- **Cobrand managers** — cobrand admins of cobrands that manage the property; they need to know who they are managing for, but cannot change ownership.

Tenants, household members, and agents do not see ownership records.

**create** is open to any authenticated user. The backend `PropertyOwnerApplicationService` enforces the bootstrap rule (the first owner of a brand-new property is created in the same transaction as the property) and the existing-owner rule (only an existing owner of a property may add another `propertyOwners` record to it). Access rules alone cannot express the bootstrap case without a chicken-and-egg failure on the first record.

**update / delete** — restricted to current person and cobrand owners only:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
```

Cobrand managers may not modify or remove ownership records. The application service must additionally prevent deletion of the last remaining owner record for a property.

#### `cobrands` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | `@request.auth.id != ""` |
| update | *(see below)* |
| delete | *(see below — same as update)* |

**list / view** — a cobrand is visible only to its own admins and its agents:

```
cobrandAdmins.user.id = @request.auth.id
|| agents.person.user.id = @request.auth.id
```

**create** is open to any authenticated user so that a person can form a new cobrand and bootstrap themselves as its first admin.

**update / delete** — restricted to cobrand admins only; agents may not modify the cobrand record:

```
cobrandAdmins.user.id = @request.auth.id
```

#### `cobrandAdmins` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | `@request.auth.id != ""` |
| update | *(see below)* |
| delete | *(see below — differs from update)* |

**list / view** — visible to existing admins of the same cobrand and to agents of that cobrand:

```
cobrand.cobrandAdmins.user.id = @request.auth.id
|| cobrand.agents.person.user.id = @request.auth.id
```

Agents need to know who their principals are; they can already see the cobrand record itself.

**create** is open to any authenticated user. `CobrandApplicationService` enforces the bootstrap rule (the first admin is created atomically with the cobrand) and the existing-admin rule (only an existing admin may add further admins to an already-administered cobrand). The access rule alone cannot express the bootstrap case without a chicken-and-egg failure on the first record.

**update** — restricted to existing cobrand admins:

```
cobrand.cobrandAdmins.user.id = @request.auth.id
```

**delete** — existing cobrand admins may remove any admin record; additionally, the linked user may remove themselves (resignation):

```
cobrand.cobrandAdmins.user.id = @request.auth.id
|| user.id = @request.auth.id
```

#### `personPropertyOwners` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | *(see below)* |
| update | *(see below — same as create)* |
| delete | *(see below — same as create)* |

**list / view** — visible to owners and managers of the linked property (same three groups as `propertyOwners`, routed through `propertyOwner.property`):

```
propertyOwner.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| propertyOwner.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| propertyOwner.property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
```

Tenants, household members, and agents do not see ownership sub-records.

**create / update / delete** — restricted to current owners only; managers may not modify who is linked as an owner:

```
propertyOwner.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| propertyOwner.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
```

#### `cobrandPropertyOwners` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | *(see below)* |
| update | *(see below — same as create)* |
| delete | *(see below — differs from create)* |

**list / view** — same three groups as `personPropertyOwners`:

```
propertyOwner.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| propertyOwner.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| propertyOwner.property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
```

**create / update** — restricted to current property owners only:

```
propertyOwner.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| propertyOwner.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
```

**delete** — property owners may remove a cobrand owner; additionally, the cobrand being removed may resign its own ownership:

```
propertyOwner.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| propertyOwner.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| cobrand.cobrandAdmins.user.id = @request.auth.id
```

#### `cobrandPropertyManagers` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | *(see below)* |
| update | *(see below — same as create)* |
| delete | *(see below — differs from create)* |

**list / view** — property owners can see who manages their property; the managing cobrand's own admins can see their own management record:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| cobrand.cobrandAdmins.user.id = @request.auth.id
```

The three groups:
- **Person owners** — owners of the linked property.
- **Cobrand owners** — cobrand admins of cobrands that own the linked property.
- **Managing cobrand admins** — admins of the cobrand that holds this management record; they need to see the properties they manage.

**create / update** — only property owners may appoint or change managers:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
```

**delete** — property owners may revoke management; the managing cobrand's admins may resign:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| cobrand.cobrandAdmins.user.id = @request.auth.id
```

#### `agents` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | *(see below)* |
| update | *(see below — same as create)* |
| delete | *(see below — differs from create)* |

**list / view** — visible to cobrand admins of the linked cobrand, the agent themselves, and owners/managers of any property this agent is assigned to:

```
cobrand.cobrandAdmins.user.id = @request.auth.id
|| person.user.id = @request.auth.id
|| propertyAgents.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| propertyAgents.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| propertyAgents.property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
```

**create / update** — only cobrand admins may register a person as an agent of their cobrand or change that record:

```
cobrand.cobrandAdmins.user.id = @request.auth.id
```

**delete** — cobrand admins may remove any agent; the agent themselves may resign:

```
cobrand.cobrandAdmins.user.id = @request.auth.id
|| person.user.id = @request.auth.id
```

#### `propertyAgents` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | *(see below)* |
| update | *(see below — same as create)* |
| delete | *(see below — differs from create)* |

**list / view** — visible to property owners/managers and to the agent and their cobrand admins:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
|| agent.person.user.id = @request.auth.id
|| agent.cobrand.cobrandAdmins.user.id = @request.auth.id
```

**create / update** — either side may initiate or change the assignment: property owners, cobrand owners, cobrand property managers, or the agent's cobrand admins:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
|| agent.cobrand.cobrandAdmins.user.id = @request.auth.id
```

**delete** — same as create, plus the agent themselves (resignation):

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
|| agent.cobrand.cobrandAdmins.user.id = @request.auth.id
|| agent.person.user.id = @request.auth.id
```

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

properties ── propertyHistories
           ── tenantHistories ── persons
           ── propertyItemHistories ── items
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
- `PropertyHistoryRepository`
- `TenantHistoryRepository`
- `PropertyItemHistoryRepository`

### 4.2 Modified repositories

- `PersonRepository` — update queries to remove `type`/`property` fields; add `user`, `DOB`.
- `PropertyRepository` — remove `owners`/`managers` field references.
- Delete `DeviceRepository`, `PersonDeviceRepository`.

### 4.3 New services (`internal/services/`)

One service per entity group, containing business logic and validation. Application services call these; services call repositories directly.

- `ItemService` — item validation, entry device state transitions, defunct reason enforcement
- `CobrandService` — cobrand validation, admin uniqueness
- `AgentService` — agent registration rules, duplicate assignment prevention
- `PropertyOwnerService` — ownership validation, duplicate person/cobrand owner prevention
- `PersonService` — person validation
- `PropertyService` — property validation, ownership chain resolution

### 4.4 New application services (`internal/application/`)

One application service per entity group, coordinating service calls without containing business logic. Each is injected with its required services via the `dig` container.

- `ItemApplicationService` — orchestrates item CRUD, entry device designation, and person/property item associations
- `CobrandApplicationService` — orchestrates cobrand CRUD, admin assignment, and property manager assignment
- `AgentApplicationService` — orchestrates agent registration and property assignment
- `PropertyOwnerApplicationService` — orchestrates property owner creation and person/cobrand owner linkage
- `PersonApplicationService` — orchestrates person CRUD and user account linking
- `PropertyApplicationService` — orchestrates property CRUD

Hook handlers in `cmd/keybook.go` call application services; application services call service-layer functions; services call repositories.

### 4.5 History services and hooks

Three history services record property-scoped change events. Remove `PersonHistoryServices`, `PropertyHistoryServices`, `DeviceHistoryServices`, and `PersonDeviceHistoryServices` after migration.

- `PropertyHistoryService` — called from `OnModelAfterCreate` and `OnModelBeforeUpdate` hooks on `properties`; writes a `propertyHistories` record with the full property state in `snapshot`.
- `TenantHistoryService` — called from `OnModelAfterCreate` and `OnModelAfterDelete` hooks on `tenants`; writes a `tenantHistories` record with action `Added` or `Removed`, capturing full `personSnapshot` and `propertySnapshot` from the live records before any deletion occurs.
- `PropertyItemHistoryService` — called from `OnModelAfterCreate` and `OnModelAfterDelete` hooks on `propertyItems`; writes a `propertyItemHistories` record with action `Added` or `Removed`, capturing full `itemSnapshot` and `propertySnapshot` from the live records before any deletion occurs.

Add `PropertyHistoryRepository`, `TenantHistoryRepository`, and `PropertyItemHistoryRepository` to `internal/repositories/` to support these services.

### 4.6 DTOs (`internal/dtos/`)

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

- `person/` — remove type/property handling; add user-link, DOB. Role derivation (owner/tenant/agent/household) is delegated to `PersonService` — modules do not query relation collections directly.
- `property/` — remove direct owners/managers. Ownership chain resolution is delegated to `PropertyService` — modules do not query ownership collections directly.

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

1. Write a PocketBase migration in `backend/migrations/` that defines all collections as specified in section 2.1, sets unique constraints, and configures access rules using the new ownership chain. Delete `database/pb_schema.json`.
2. On the OpenBSD server: stop the binary, delete `pb_data/`, build and run the new binary (migrations run automatically on `serve`), restart.
3. Migrate existing data (run as a one-shot script):
   - For each `devices` record → create `items` + `entryDevices` record.
   - For each `personDevices` record → create `personItems` record.
   - For each `persons` record with type `Owner` → create `propertyOwners` + `personPropertyOwners` record.
   - For each `persons` record with type `Tenant` → create `tenants` record.
   - For each `persons` record with type `Agent` → create `agents` + `propertyAgents` record.
   - For each `persons` record with type `Household` → create `households` record.
   - For each `properties` record → create `propertyOwners` records from the `owners` relation.
4. Update all PocketBase access rules to use the new ownership chain.
