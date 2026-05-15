# Design — full-collections

## 1. Overview

This change restructures the data model from a flat, property-centric design to a normalised multi-role design. Persons become first-class entities (not property-scoped). Ownership, tenancy, and management are expressed through dedicated relation collections rather than embedded fields. A new item registry supports physical household items and entry devices. Firm (cobrand) ownership and agent management are introduced as first-class concepts.

The old database is dropped and replaced from scratch. Columns from old collections are carried forward onto their counterpart collections in the new schema where they are not already covered. `database/pb_schema.json` is deleted — schema is defined and managed exclusively via PocketBase migrations in `backend/migrations/`.

## 2. Schema changes

### 2.1 Complete collection definitions

All columns listed. Columns marked *(carried over)* come from the old schema. History collections follow the existing pattern (snapshot JSON + `statedDateTime` per entity) and are retained for `persons` and `properties` only — relation-only tables rely on PocketBase's built-in `created`/`updated` timestamps.

---

**`users`** *(auth)*
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

### 2.5 Access rule implications

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

### 4.1 New repositories

One repository per new collection following the existing pattern in `backend/internal/repositories/`:

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
- Remove `DeviceRepository`, `PersonDeviceRepository` after migration.

### 4.3 Services and hooks

No new history services are required for the relation collections. The existing `PersonHistoryServices` and `PropertyHistoryServices` hooks continue to apply. Remove `DeviceHistoryServices` and `PersonDeviceHistoryServices` after migration.

### 4.4 DTOs

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

### 5.2 New modules (`src/lib/modules/`)

- `item/` — `ItemListModule`, `ItemDetailModule`
- `cobrand/` — `CobrandListModule`, `CobrandDetailModule`
- `agent/` — `AgentListModule`, `AgentDetailModule`

### 5.3 Updated modules

- `person/` — remove type/property handling; add user-link, DOB, role derivation via relation collections.
- `property/` — remove direct owners/managers; derive ownership from new collections.

### 5.4 New contexts

- `ItemContext` (following existing `DeviceContext` pattern)
- `CobrandContext`
- `AgentContext`

## 6. Migration path

1. Add all new collections to `pb_schema.json`.
2. Migrate existing data:
   - For each `devices` record → create `items` + `entryDevices` record.
   - For each `personDevices` record → create `personItems` record.
   - For each `persons` record with type `Owner` → create `propertyOwners` + `personPropertyOwners` record.
   - For each `persons` record with type `Tenant` → create `tenants` record.
   - For each `persons` record with type `Agent` → create `agents` + `propertyAgents` record.
   - For each `persons` record with type `Household` → create `households` record.
   - For each `properties` record → create `propertyOwners` records from the `owners` relation.
3. Remove old fields (`persons.type`, `persons.property`, `properties.owners`, `properties.managers`).
4. Remove old collections (`devices`, `personDevices`, `deviceHistory`, `personDeviceHistory`).
5. Update all PocketBase access rules to use the new ownership chain.
