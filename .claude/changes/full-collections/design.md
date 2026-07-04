# Design — full-collections

## 1. Overview

This change restructures the data model from a flat, property-centric design to a normalised multi-role design. Persons become first-class entities (not property-scoped). Ownership, tenancy, and management are expressed through dedicated relation collections rather than embedded fields. A new item registry supports physical household items and entry devices. Firm (cobrand) ownership and agent management are introduced as first-class concepts.

The old database is dropped and replaced from scratch. Columns from old collections are carried forward onto their counterpart collections in the new schema where they are not already covered. The new schema is defined and managed exclusively via PocketBase migrations in `backend/migrations/`. `database/pb_schema.json` is deleted.

All new backend and frontend code follows the layered architecture defined in `CLAUDE.md`: API → Application → Service → Repository → Store. This change introduces `internal/application/` (backend) and `src/lib/services/`, `src/lib/repositories/` (frontend) as new layers.

Registration creates only the user account. Post-login routing then sends the user to their property list (linked person), their cobrand list (linked cobrandAdmins record), or an account-setup choice page offering either path (neither) — section 5.7. Cobrand admins are additionally gated by an `approved` flag (sections 2.1, 2.4) that only a KeyBook staff member can set; an unapproved admin can view their own pending cobrand but cannot create, update, or delete cobrand-scoped records.

## 2. Schema changes

### 2.1 Complete collection definitions

All columns listed. Columns marked *(carried over)* come from the old schema. Relation-only tables rely on PocketBase's built-in `created`/`updated` timestamps.

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
| `approved` | bool, default `false` | gates cobrand-management access (section 2.4); settable only by a KeyBook staff member via the PocketBase superuser dashboard — no app-facing update path exists |

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

### 2.2 Collections not carried forward

The following old collections are dropped with the database and have no equivalent in the new schema:

- `devices` — superseded by `items` + `entryDevices`
- `personDevices` — superseded by `personItems`
- `deviceHistory` — dropped (entry devices tracked via `entryDevices` timestamps)
- `personDeviceHistory` — dropped
- `personHistory` — dropped; person changes are not tracked in the new schema
- `propertyHistory` — dropped; property-scoped history is not tracked in the new schema

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

List/view visibility here is unaffected by cobrand approval status — an unapproved cobrand admin/manager can still see tenants and item-holders at properties their (pending) cobrand touches, same as any other read.

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
|| (propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
```

Cobrand managers, tenants, household members, and agents can view but not modify the property record itself. An unapproved cobrand owner is likewise blocked (section 2.4).

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
|| (property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
```

Cobrand managers may not modify or remove ownership records, nor may an unapproved cobrand owner. The application service must additionally prevent deletion of the last remaining owner record for a property.

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

Visibility is unaffected by approval status — an unapproved admin still needs to see their own pending cobrand.

**create** is open to any authenticated user so that a person can form a new cobrand and bootstrap themselves as its first (unapproved) admin.

**update / delete** — restricted to approved cobrand admins only; agents may not modify the cobrand record:

```
cobrandAdmins.user.id = @request.auth.id && cobrandAdmins.approved = true
```

#### `cobrandAdmins` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | `@request.auth.id != ""` |
| update | `null` (superusers only) |
| delete | *(see below)* |

**list / view** — visible to existing admins of the same cobrand and to agents of that cobrand:

```
cobrand.cobrandAdmins.user.id = @request.auth.id
|| cobrand.agents.person.user.id = @request.auth.id
```

Agents need to know who their principals are; they can already see the cobrand record itself. Visibility is unaffected by approval status.

**create** is open to any authenticated user. `CobrandApplicationService` enforces the bootstrap rule (the first admin is created atomically with the cobrand, regardless of approval, since none can exist yet) and the existing-admin rule (only an existing admin may add further admins to an already-administered cobrand, and that inviting admin's own `approved` must be `true` — section 2.4). The access rule alone cannot express the bootstrap case without a chicken-and-egg failure on the first record.

**update** — locked to PocketBase superusers (`null` rule). No authenticated app user, including existing admins, may update this collection through the API. This is the only mechanism for setting `approved` (section 2.4), and it is deliberately kept outside the app's write path — "contact KeyBook" means a staff member flips the field directly in the PocketBase dashboard.

**delete** — existing *approved* cobrand admins may remove any admin record; additionally, the linked user may remove themselves (resignation) regardless of approval:

```
(cobrand.cobrandAdmins.user.id = @request.auth.id && cobrand.cobrandAdmins.approved = true)
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

**create / update / delete** — restricted to current owners only; managers may not modify who is linked as an owner, nor may an unapproved cobrand owner:

```
propertyOwner.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (propertyOwner.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && propertyOwner.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
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

**create / update** — restricted to current property owners only; an unapproved cobrand owner may not act here:

```
propertyOwner.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (propertyOwner.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && propertyOwner.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
```

**delete** — property owners may remove a cobrand owner (an unapproved cobrand owner may not remove a *different* cobrand's ownership); additionally, the cobrand being removed may resign its own ownership regardless of approval:

```
propertyOwner.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (propertyOwner.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && propertyOwner.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
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

**create / update** — only property owners may appoint or change managers; an unapproved cobrand owner may not:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
```

**delete** — property owners may revoke management (an unapproved cobrand owner may not); the managing cobrand's admins may resign regardless of approval:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
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

**create / update** — only *approved* cobrand admins may register a person as an agent of their cobrand or change that record:

```
cobrand.cobrandAdmins.user.id = @request.auth.id && cobrand.cobrandAdmins.approved = true
```

**delete** — approved cobrand admins may remove any agent; the agent themselves may resign regardless of the cobrand's approval status:

```
(cobrand.cobrandAdmins.user.id = @request.auth.id && cobrand.cobrandAdmins.approved = true)
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

**create / update** — either side may initiate or change the assignment: property owners, cobrand owners, cobrand property managers, or the agent's cobrand admins — each cobrand-based party must be approved:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
|| (property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id && property.cobrandPropertyManagers.cobrand.cobrandAdmins.approved = true)
|| (agent.cobrand.cobrandAdmins.user.id = @request.auth.id && agent.cobrand.cobrandAdmins.approved = true)
```

**delete** — same approval-gated parties as create/update, plus the agent themselves (resignation, ungated):

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
|| (property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id && property.cobrandPropertyManagers.cobrand.cobrandAdmins.approved = true)
|| (agent.cobrand.cobrandAdmins.user.id = @request.auth.id && agent.cobrand.cobrandAdmins.approved = true)
|| agent.person.user.id = @request.auth.id
```

#### `items` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | `@request.auth.id != ""` |
| update | *(see below)* |
| delete | *(see below — same as update)* |

**list / view** — visible to person-owners, property owners/managers, and any resident (tenant, household member, or agent) of a property the item is assigned to:

```
personItems.person.user.id = @request.auth.id
|| propertyItems.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| propertyItems.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| propertyItems.property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
|| propertyItems.property.tenants.person.user.id = @request.auth.id
|| propertyItems.property.households.person.user.id = @request.auth.id
|| propertyItems.property.propertyAgents.agent.person.user.id = @request.auth.id
```

**create** is open to any authenticated user so a person can register a new item before any association records exist. The `PersonApplicationService` atomically creates the `personItems` record in the same request.

**update / delete** — restricted to person-owners and property owners/managers; residents may not modify item records, nor may an unapproved cobrand owner or manager:

```
personItems.person.user.id = @request.auth.id
|| propertyItems.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (propertyItems.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && propertyItems.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
|| (propertyItems.property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id && propertyItems.property.cobrandPropertyManagers.cobrand.cobrandAdmins.approved = true)
```

#### `entryDevices` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | *(see below)* |
| update | *(see below — same as create)* |
| delete | *(see below — same as create)* |

**list / view** — same audiences as the linked item, since residents need to see entry devices for properties they occupy:

```
item.personItems.person.user.id = @request.auth.id
|| item.propertyItems.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| item.propertyItems.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| item.propertyItems.property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
|| item.propertyItems.property.tenants.person.user.id = @request.auth.id
|| item.propertyItems.property.households.person.user.id = @request.auth.id
|| item.propertyItems.property.propertyAgents.agent.person.user.id = @request.auth.id
```

**create / update / delete** — restricted to person-owners and property owners/managers; residents may not designate, modify, or remove entry device records, nor may an unapproved cobrand owner or manager:

```
item.personItems.person.user.id = @request.auth.id
|| item.propertyItems.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (item.propertyItems.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && item.propertyItems.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
|| (item.propertyItems.property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id && item.propertyItems.property.cobrandPropertyManagers.cobrand.cobrandAdmins.approved = true)
```

#### `propertyItems` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | *(see below)* |
| update | *(see below — same as create)* |
| delete | *(see below — same as create)* |

**list / view** — visible to property owners/managers, the item's person-owners, and all residents at that property:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
|| item.personItems.person.user.id = @request.auth.id
|| property.tenants.person.user.id = @request.auth.id
|| property.households.person.user.id = @request.auth.id
|| property.propertyAgents.agent.person.user.id = @request.auth.id
```

**create / update / delete** — either side may place or remove an item: property owners/managers or the item's person-owners; an unapproved cobrand owner or manager may not:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
|| (property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id && property.cobrandPropertyManagers.cobrand.cobrandAdmins.approved = true)
|| item.personItems.person.user.id = @request.auth.id
```

#### `personItems` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | `@request.auth.id != ""` |
| update | *(see below)* |
| delete | *(see below — same as update)* |

**list / view** — visible to the person themselves and property owners/managers of properties where the item is located:

```
person.user.id = @request.auth.id
|| item.propertyItems.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| item.propertyItems.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| item.propertyItems.property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
```

**create** is open to any authenticated user. On first creation of a `personItems` record the item has no property associations yet, so a property-owner check cannot be used; the `ItemApplicationService` enforces who may claim ownership.

**update / delete** — the person themselves or property owners/managers of properties holding that item; an unapproved cobrand owner or manager may not:

```
person.user.id = @request.auth.id
|| item.propertyItems.property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (item.propertyItems.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && item.propertyItems.property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
|| (item.propertyItems.property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id && item.propertyItems.property.cobrandPropertyManagers.cobrand.cobrandAdmins.approved = true)
```

#### `households` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | *(see below)* |
| update | *(see below — same as create)* |
| delete | *(see below — differs from create)* |

**list / view** — visible to property owners/managers, any household member of the same property, and any tenant of the same property:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
|| property.households.person.user.id = @request.auth.id
|| property.tenants.person.user.id = @request.auth.id
```

`property.households.person.user.id = @request.auth.id` covers the member themselves (since their own record is included in the back-reference traversal) as well as co-members of the same household.

**create / update** — only property owners and managers may add or change household members; an unapproved cobrand owner or manager may not:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
|| (property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id && property.cobrandPropertyManagers.cobrand.cobrandAdmins.approved = true)
```

**delete** — owners and managers may remove any member (an unapproved cobrand owner or manager may not); the member themselves may leave regardless:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
|| (property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id && property.cobrandPropertyManagers.cobrand.cobrandAdmins.approved = true)
|| person.user.id = @request.auth.id
```

#### `tenants` access rules

| Operation | Rule |
|---|---|
| list | *(see below)* |
| view | *(see below — same as list)* |
| create | *(see below)* |
| update | *(see below — same as create)* |
| delete | *(see below — differs from create)* |

**list / view** — visible to property owners/managers, any tenant of the same property, and any household member of the same property:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id
|| property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id
|| property.tenants.person.user.id = @request.auth.id
|| property.households.person.user.id = @request.auth.id
```

**create / update** — only property owners and managers may add or change tenants; an unapproved cobrand owner or manager may not:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
|| (property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id && property.cobrandPropertyManagers.cobrand.cobrandAdmins.approved = true)
```

**delete** — owners and managers may remove any tenant (an unapproved cobrand owner or manager may not); the tenant themselves may vacate regardless:

```
property.propertyOwners.personPropertyOwners.person.user.id = @request.auth.id
|| (property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.user.id = @request.auth.id && property.propertyOwners.cobrandPropertyOwners.cobrand.cobrandAdmins.approved = true)
|| (property.cobrandPropertyManagers.cobrand.cobrandAdmins.user.id = @request.auth.id && property.cobrandPropertyManagers.cobrand.cobrandAdmins.approved = true)
|| person.user.id = @request.auth.id
```

### 2.4 Cobrand approval gate

`cobrandAdmins.approved` defaults to `false` and gates every write action a cobrand admin takes *as* their cobrand — creating or updating the cobrand itself, appointing agents or property managers, claiming property or item ownership, and so on. It does not gate:

- **Visibility.** List/view rules are unchanged throughout section 2.3 — an unapproved admin can still see their own cobrand, its (empty) admin/agent/property lists, and confirm the pending state.
- **Resignation.** Any clause where a cobrand admin, agent, or manager removes *their own* association (`user.id = @request.auth.id`, `person.user.id = @request.auth.id`, or a managing/owning cobrand backing out of a role it holds on someone else's property) is exempt — approval status should never trap someone into a role they want to leave.
- **Bootstrap.** The first `cobrandAdmins` record for a brand-new cobrand is created before any admin can be approved, so it is exempt by necessity.

Every other place `cobrandAdmins.user.id = @request.auth.id` (at any relation depth) authorizes a create, update, or delete elsewhere in section 2.3 is rewritten as `(...cobrandAdmins.user.id = @request.auth.id && ...cobrandAdmins.approved = true)`, using the identical relation-path prefix for both conditions so PocketBase resolves them against the same joined `cobrandAdmins` row rather than any co-admin's row happening to satisfy one half.

`approved` can only be changed by a KeyBook staff member with PocketBase superuser access — the `cobrandAdmins` update rule is locked (`null`) rather than left open to existing admins, since that's the only field on the record worth changing after creation and it must not be self-service. Approving a cobrand is therefore an out-of-band support action, not an app feature.

Existing `cobrandAdmins` rows created before this gate was introduced are backfilled to `approved = true` in the same migration that adds the column (section 8), so already-established admins are not retroactively locked out.

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
- `CobrandApplicationService` — orchestrates cobrand CRUD, admin assignment, and property manager assignment; enforces that an admin inviting another admin to an already-administered cobrand is themselves approved (section 2.4)
- `AgentApplicationService` — orchestrates agent registration and property assignment
- `PropertyOwnerApplicationService` — orchestrates property owner creation and person/cobrand owner linkage
- `PersonApplicationService` — orchestrates person CRUD and user account linking
- `PropertyApplicationService` — orchestrates property CRUD

Hook handlers in `cmd/keybook.go` call application services; application services call service-layer functions; services call repositories.

### 4.5 Legacy hook removal

Remove `PersonHistoryServices`, `PropertyHistoryServices`, `DeviceHistoryServices`, and `PersonDeviceHistoryServices` hooks and their DI registrations from `cmd/keybook.go` after migration.

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
| `/user/setup/` | Account setup choice (onboarding) — pick person or cobrand admin |
| `/user/persons/setup/` | Person setup (onboarding) — create the logged-in user's own linked person |

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
- `CobrandService` — cobrand validation, admin uniqueness enforcement, resolving a user's admin record
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

### 5.7 User onboarding flow

Registration creates only the user account. Post-login routing at `/user` then branches three ways by what's already linked to the account, and an account-setup page offers the choice when neither exists. All of this is frontend-only — no backend or schema changes beyond the `approved` field itself (sections 2.1, 2.4).

**Registration** (`/auth/register`):

```
RegisterForm (email, password, passwordConfirm)
  → RegisterModule.callApi()
      1. users.create({ email, password, passwordConfirm })   — SDK (auth collection)
      2. users.authWithPassword(email, password)              — SDK
  → page stores the returned auth cookie → goto /user → post-login routing
```

`RegisterModule` is unchanged from its original shape (email/password/passwordConfirm only), except that it must stop sending `name` to `users.create()` — that field was removed from `users` (section 2.1), so it is silently dropped today. `RegisterForm` drops its name input to match. No person or cobrand is created at registration time; onboarding happens after login, driven by routing.

**Post-login routing** (`/user/+page.ts`):

The existing unconditional redirect to `/user/properties/list` becomes a three-way branch. The load function resolves `PersonService.getPersonByUserId(authUserId)` first and only falls through to the cobrand-admin check if no person is linked:

| Linked person | Linked cobrandAdmins | Redirect |
|---|---|---|
| found | — | `/user/properties/list` |
| none | found | `/user/cobrands/` |
| none | none | `/user/setup` |

Both the login page and the register page navigate to `/user`, so this one decision point covers both entry paths as well as direct navigation. `CobrandService` gains `getAdminRecordForUserId(userId): Promise<ICobrandAdminModel | null>`, wrapping a new `CobrandAdminRepository.getByUserId(userId)` — the same shape as `PersonRepository.getByUserId` / `PersonService.getPersonByUserId`.

**Account setup choice** (`/user/setup/`):

A page offering two buttons: "Set up my person profile" (→ `/user/persons/setup/`) and "Set up my company" (→ `/user/cobrands/add/`, the existing Phase 8 cobrand-creation route — no new cobrand UI is introduced). If the user already resolves to a person or a cobrand admin (e.g. reached via the back button after already onboarding), the page redirects per the table above instead of showing the choice.

New module: `AccountSetupModule` (`src/lib/modules/user/`) — resolves the same two lookups as the `/user/+page.ts` load function to decide whether to show the choice or redirect away.

**Person setup page** (`/user/persons/setup/`):

Unchanged from the prior design: a form collecting name and date of birth for the user's own person. On submit it calls `PersonService.createPerson(name, dob, authUserId)` and navigates to `/user/properties/list`. If the user already has a linked person, the page redirects to the property list. It remains deliberately separate from `/user/persons/add/`, which creates unlinked persons (owners registering tenants who have no account).

**Cobrand admin setup** (existing `/user/cobrands/add/`):

No new route or component — the choice page links directly to the flow already built in Phase 8 (`CobrandAddModule`: create the cobrand, then call `CobrandService.addAdmin` to link the current user as its first admin). The only change is presentational: since `approved` defaults to `false`, the cobrand detail page must show a pending-approval notice ("Contact KeyBook to activate management access for this cobrand") whenever the viewing user's own admin record on that cobrand has `approved = false`. Because the underlying access rules already permit list/view regardless of approval, this is a display-only addition — no new query is needed beyond what the detail page already loads.

**Layer changes:**

- `CobrandAdminRepository.getByUserId(userId)` — new; mirrors `PersonRepository.getByUserId`.
- `CobrandService.getAdminRecordForUserId(userId)` — new; wraps the repository method.
- `PersonRepository.create` / `PersonService.createPerson` — unchanged from the prior design (still accept an optional `userId` to set the link); exercised by `/user/persons/setup/`, no longer by registration.

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

### Contract test structure

Contract tests are Go tests that start a real PocketBase HTTP server (`t.TempDir()` data directory), bootstrap a shared fixture set via the admin API, then make direct HTTP requests authenticated as each user type to verify correct status codes.

**Standard fixture set** — created once per test file:

| Fixture | Type | Key relationships |
|---|---|---|
| `userOwner` | users | linked to `personOwner` |
| `userCobrandAdmin` | users | admin of `cobrand`; cobrand co-owns `property` via `cobrandPropertyOwners` |
| `userManager` | users | cobrand manager of `property` via `cobrandPropertyManagers` |
| `userTenant` | users | linked to `personTenant` |
| `userHousehold` | users | linked to `personHousehold` |
| `userAgent` | users | linked to `personAgent`; assigned to `property` via `propertyAgents` |
| `userUnrelated` | users | no relationship to any entity |
| `userUnapprovedAdmin` | users | a second admin of `cobrand`, not yet approved |
| `personOwner` | persons | user = `userOwner` |
| `personTenant` | persons | user = `userTenant` |
| `personHousehold` | persons | user = `userHousehold` |
| `personAgent` | persons | user = `userAgent` |
| `cobrand` | cobrands | |
| `cobrandAdmin` | cobrandAdmins | user = `userCobrandAdmin`, cobrand = `cobrand`, `approved = true` |
| `unapprovedCobrandAdmin` | cobrandAdmins | user = `userUnapprovedAdmin`, cobrand = `cobrand`, `approved = false` |
| `property` | properties | |
| `propertyOwner` | propertyOwners | property = `property` |
| `personPropertyOwner` | personPropertyOwners | person = `personOwner`, propertyOwner = `propertyOwner` |
| `cobrandPropertyOwner` | cobrandPropertyOwners | cobrand = `cobrand`, propertyOwner = `propertyOwner` |
| `cobrandPropertyManager` | cobrandPropertyManagers | cobrand = `cobrand`, property = `property` |
| `tenant` | tenants | person = `personTenant`, property = `property` |
| `household` | households | person = `personHousehold`, property = `property` |
| `agent` | agents | person = `personAgent`, cobrand = `cobrand` |
| `propertyAgent` | propertyAgents | agent = `agent`, property = `property` |
| `item` | items | |
| `personItem` | personItems | person = `personOwner`, item = `item` |
| `propertyItem` | propertyItems | item = `item`, property = `property` |
| `entryDevice` | entryDevices | item = `item` |

**Test pattern** — for each collection, assert:

1. Unauthenticated requests to all five operations (list, view, create, update, delete) return 403.
2. Each user type in the authorized group for that operation returns 200 (list/view), 201 (create), or 204 (update/delete).
3. Each user type that should be blocked returns 403.
4. Special cases (bootstrap, last-owner guard, self-removal, resignation, cobrand approval gate) are tested at whichever layer enforces them — application-service-level behaviours verify the correct HTTP status and error body returned by the application service; the approval gate is a raw access-rule behaviour and is asserted directly. For the approval gate specifically: `userUnapprovedAdmin` gets 403 on every create/update/delete gated in section 2.4 despite being a genuine admin of `cobrand`; `userCobrandAdmin` (approved) continues to succeed on those same operations (regression check against the original 1.1–1.3 expectations); both get 200 on list/view of `cobrand`; `userUnapprovedAdmin` can still delete their own `unapprovedCobrandAdmin` record (resignation); no authenticated user, including `userCobrandAdmin`, can update any `cobrandAdmins` record.

Expected outcomes for each user type × collection × operation are derived directly from the EARS requirements in `requirements.md`.

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
5. A later migration (section 2.4) adds `cobrandAdmins.approved` (default `false`) and backfills existing rows to `approved = true` so already-established admins are not retroactively locked out; only rows created after this migration default to `false`. It also rewrites the access rules listed in section 2.3 to include the approval gate and locks the `cobrandAdmins` update rule to superusers only.
