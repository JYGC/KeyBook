# Tasks — full-collections

## Phase 1 — Schema

- [x] **1.1** Write contract tests for `users`, `persons`, and `properties` using the standard fixture set from `design.md` section 7. For `persons`: verify own user, tenant-property-owner, and item-property-owner groups get 200 on list/view; unrelated user gets 403; any authenticated user can create (201); own user can update/delete (200); non-linked user gets 403 on update/delete. For `properties`: verify all six visibility groups get 200 on list/view; unrelated gets 403; any authenticated can create (201); only owners (person or cobrand) can update/delete (200); other groups get 403 on update/delete.
- [x] **1.2** Write contract tests for `propertyOwners`, `personPropertyOwners`, `cobrandPropertyOwners`, and `cobrandPropertyManagers`. For `propertyOwners`: verify person owner, cobrand owner, and cobrand manager get 200 on list/view; tenant and unrelated get 403; any authenticated can create a first record (bootstrap, 201); non-owner cannot create a second record; manager cannot update/delete (403); last-owner deletion is rejected by the application service. For `personPropertyOwners` and `cobrandPropertyOwners`: verify same three-group visibility; only owners can create/update/delete (403 for managers); for `cobrandPropertyOwners`, verify cobrand admin of the linked cobrand can delete their own record. For `cobrandPropertyManagers`: verify property owners and managing cobrand admin can list/view; only property owners can create/update (403 for managing cobrand admin); both property owners and managing cobrand admin can delete (200).
- [x] **1.3** Write contract tests for `cobrands`, `cobrandAdmins`, `agents`, `propertyAgents`, `households`, `tenants`, `items`, `entryDevices`, `propertyItems`, and `personItems`. Key cases: cobrand admin and agent can see cobrands (200), unrelated cannot (403); cobrand admin bootstrap (first cobrandAdmin created atomically with cobrand, non-admin cannot add admins to existing cobrand); cobrand admin can self-remove from cobrandAdmins; cobrand admin and agent can each delete the agents record; agent can delete own propertyAgents record (resignation); tenant can see household records and household member can see tenant records (cross-visibility, 200); member and tenant can each self-delete; residents (tenant, household member, agent) can list/view items and entryDevices (200) but cannot create/update/delete entryDevices (403); either the item's person-owner or the property's owner/manager can create/delete propertyItems and personItems (201/204).
- [x] **1.4** Write a PocketBase migration in `backend/migrations/` that creates all collections as specified in `design.md` section 2.1, with all unique constraints, carry-over columns (`persons.profileImage`, `entryDevices.identifier`, `entryDevices.defunctReason`), and removal of default `name` and `avatar` fields from `users`. Delete `database/pb_schema.json`. No access rules yet — set all rules to admin-only (`null`) as a safe placeholder.
- [x] **1.5** In the same migration, add access rules for `users`, `persons`, and `properties` per `design.md` section 2.3.
- [x] **1.6** In the same migration, add access rules for the ownership chain collections — `propertyOwners`, `personPropertyOwners`, `cobrandPropertyOwners`, `cobrandPropertyManagers` — per `design.md` section 2.3. Note the bootstrap exception for `propertyOwners` create (enforcement delegated to `PropertyOwnerApplicationService`) and the last-owner-deletion guard (application service only).
- [x] **1.7** In the same migration, add access rules for the remaining collections — `cobrands`, `cobrandAdmins`, `agents`, `propertyAgents`, `households`, `tenants`, `items`, `entryDevices`, `propertyItems`, `personItems` — per `design.md` section 2.3.
- [x] **1.8** Confirm all contract tests from 1.1–1.3 pass against the completed migration.
- [x] **1.9** On the OpenBSD server: stop the binary, delete `pb_data/`, build and run the new binary (migrations run automatically on `serve`). Verify all collections exist with correct fields and rules.

## Phase 2 — Backend: DTOs

- [x] **2.1** Write unit tests for all new and updated DTOs covering field mapping and required-field validation.
- [x] **2.2** Add DTOs: `CobrandDto`, `CobrandAdminDto`, `CobrandPropertyManagerDto`, `CobrandPropertyOwnerDto`.
- [x] **2.3** Add DTOs: `PropertyOwnerDto`, `PersonPropertyOwnerDto`.
- [x] **2.4** Add DTOs: `AgentDto`, `PropertyAgentDto`.
- [x] **2.5** Add DTOs: `HouseholdDto`, `TenantDto`.
- [x] **2.6** Add DTOs: `ItemDto`, `PropertyItemDto`, `PersonItemDto`, `EntryDeviceDto` (include `identifier` and `defunctReason`).
- [x] **2.7** Rewrite `PersonDto`: fields are `name`, `DOB`, `user`, `profileImage`.
- [x] **2.8** Rewrite `PropertyDto`: fields are `address` only.

## Phase 3 — Backend: Repositories

- [x] **3.1** Write integration tests for each new repository using a real PocketBase HTTP server (`t.TempDir()`) — no database mocking. Cover CRUD operations and relation queries.
- [x] **3.2** Add `ItemRepository`, `PropertyItemRepository`, `PersonItemRepository`, `EntryDeviceRepository`.
- [x] **3.3** Add `PropertyOwnerRepository`, `PersonPropertyOwnerRepository`.
- [x] **3.4** Add `CobrandRepository`, `CobrandAdminRepository`, `CobrandPropertyManagerRepository`, `CobrandPropertyOwnerRepository`.
- [x] **3.5** Add `AgentRepository`, `PropertyAgentRepository`.
- [x] **3.6** Add `HouseholdRepository`, `TenantRepository`.
- [x] **3.7** Rewrite `PersonRepository`: fields are `name`, `DOB`, `user`, `profileImage`. Confirm integration tests pass.
- [x] **3.8** Rewrite `PropertyRepository`: fields are `address` only. Confirm integration tests pass.
- [x] **3.9** Delete `DeviceRepository` and `PersonDeviceRepository`.

## Phase 4 — Backend: Services

- [ ] **4.1** Write unit tests for each new service covering business logic and validation in isolation (no PocketBase dependency).
- [ ] **4.2** Add `ItemService` (item validation, entry device state transitions, defunct reason enforcement).
- [ ] **4.3** Add `CobrandService` (cobrand validation, admin uniqueness).
- [ ] **4.4** Add `AgentService` (agent registration rules, duplicate assignment prevention).
- [ ] **4.5** Add `PropertyOwnerService` (ownership validation, duplicate person/cobrand owner prevention).
- [ ] **4.6** Add `PersonService` (person validation).
- [ ] **4.7** Add `PropertyService` (property validation, ownership chain resolution).

## Phase 5 — Backend: Application layer

- [ ] **5.1** Write integration tests for each application service using a real PocketBase HTTP server (`t.TempDir()`) — cover full use-case orchestration flows including cross-service operations.
- [ ] **5.2** Add `ItemApplicationService` (orchestrates item CRUD, entry device designation, person/property item associations via `ItemService`).
- [ ] **5.3** Add `CobrandApplicationService` (orchestrates cobrand CRUD, admin assignment, property manager assignment via `CobrandService`).
- [ ] **5.4** Add `AgentApplicationService` (orchestrates agent registration and property assignment via `AgentService`).
- [ ] **5.5** Add `PropertyOwnerApplicationService` (orchestrates property owner creation, person/cobrand owner linkage via `PropertyOwnerService`).
- [ ] **5.6** Add `PersonApplicationService` (orchestrates person CRUD and user account linking via `PersonService`).
- [ ] **5.7** Add `PropertyApplicationService` (orchestrates property CRUD via `PropertyService`).

## Phase 6 — Backend: DI wiring

- [ ] **6.1** Register all new repositories, services, and application services in the `dig` container in `cmd/keybook.go`.
- [ ] **6.2** Update hook handlers in `cmd/keybook.go` to call application services rather than services or repositories directly.
- [ ] **6.3** Remove `PersonHistoryServices`, `PropertyHistoryServices`, `DeviceHistoryServices`, and `PersonDeviceHistoryServices` hooks and their registrations from `cmd/keybook.go`.
- [ ] **6.4** Remove `DeviceRepository` and `PersonDeviceRepository` from the DI container.

## Phase 7 — Frontend: Item registry

- [ ] **7.1** Write unit tests for `ItemService` business logic (entry device state transitions, defunct reason validation).
- [ ] **7.2** Write integration tests for `ItemRepository`, `EntryDeviceRepository`, `PropertyItemRepository`, `PersonItemRepository` against a real running PocketBase instance — no SDK mocking.
- [ ] **7.3** Implement `src/lib/repositories/item/` — `ItemRepository`, `EntryDeviceRepository`, `PropertyItemRepository`, `PersonItemRepository`.
- [ ] **7.4** Implement `src/lib/services/item/` — `ItemService`.
- [ ] **7.5** Write unit tests for `ItemListModule` and `ItemDetailModule`.
- [ ] **7.6** Implement `src/lib/modules/item/` — `ItemListModule`, `ItemDetailModule` (calls `ItemService`; does not call repositories or SDK directly).
- [ ] **7.7** Add `ItemContext` to the user layout.
- [ ] **7.8** Implement components: `ItemList`, `ItemEditor`, `EntryDeviceEditor` (identifier, deviceType, defunctReason).
- [ ] **7.9** Add routes: `/user/items/`, `/user/items/add/`, `/user/items/edit/`.
- [ ] **7.10** Write E2E tests for item create, edit, delete, entry device designation, and defunct marking flows.

## Phase 8 — Frontend: Cobrand management

- [ ] **8.1** Write unit tests for `CobrandService` business logic (admin uniqueness enforcement).
- [ ] **8.2** Write integration tests for cobrand repositories against a real running PocketBase instance — no SDK mocking.
- [ ] **8.3** Implement `src/lib/repositories/cobrand/` — `CobrandRepository`, `CobrandAdminRepository`, `CobrandPropertyManagerRepository`, `CobrandPropertyOwnerRepository`.
- [ ] **8.4** Implement `src/lib/services/cobrand/` — `CobrandService`.
- [ ] **8.5** Write unit tests for `CobrandListModule` and `CobrandDetailModule`.
- [ ] **8.6** Implement `src/lib/modules/cobrand/` — `CobrandListModule`, `CobrandDetailModule`.
- [ ] **8.7** Add `CobrandContext` to the user layout.
- [ ] **8.8** Implement components: `CobrandList`, `CobrandEditor`, `CobrandAdminList`, `CobrandPropertyManagerList`.
- [ ] **8.9** Add routes: `/user/cobrands/`, `/user/cobrands/add/`, `/user/cobrands/[id]/`.
- [ ] **8.10** Write E2E tests for cobrand create, admin assignment, and property manager assignment.

## Phase 9 — Frontend: Agent management

- [ ] **9.1** Write unit tests for `AgentService` business logic (registration rules, duplicate assignment prevention).
- [ ] **9.2** Write integration tests for agent repositories against a real running PocketBase instance — no SDK mocking.
- [ ] **9.3** Implement `src/lib/repositories/agent/` — `AgentRepository`, `PropertyAgentRepository`.
- [ ] **9.4** Implement `src/lib/services/agent/` — `AgentService`.
- [ ] **9.5** Write unit tests for `AgentListModule` and `AgentDetailModule`.
- [ ] **9.6** Implement `src/lib/modules/agent/` — `AgentListModule`, `AgentDetailModule`.
- [ ] **9.7** Add `AgentContext` to the user layout.
- [ ] **9.8** Implement components: `AgentList`, `AgentEditor`, `PropertyAgentList`.
- [ ] **9.9** Add routes: `/user/agents/`, `/user/agents/add/`, `/user/agents/[id]/`.
- [ ] **9.10** Write E2E tests for agent registration and property assignment.

## Phase 10 — Frontend: Persons and properties

- [ ] **10.1** Write unit tests for updated `PersonService` (role derivation — owner/tenant/agent/household via relation collections) and `PropertyService` (ownership chain traversal via propertyOwners/personPropertyOwners/cobrandPropertyOwners).
- [ ] **10.2** Write integration tests for updated person and property repositories against a real running PocketBase instance — no SDK mocking.
- [ ] **10.3** Implement/update `src/lib/repositories/person/` — `PersonRepository` and all person-relation repositories (`PersonPropertyOwnerRepository`, etc.).
- [ ] **10.4** Implement/update `src/lib/repositories/property/` — `PropertyRepository` and all property-relation repositories.
- [ ] **10.5** Implement/update `src/lib/services/person/` — `PersonService` (role derivation delegated from PersonModule; does not call SDK directly).
- [ ] **10.6** Implement/update `src/lib/services/property/` — `PropertyService` (ownership chain resolution delegated from PropertyModule; does not call SDK directly).
- [ ] **10.7** Write unit tests for updated `PersonModule` and `PropertyModule`.
- [ ] **10.8** Rewrite `PersonModule`: fields are name, DOB, profileImage, user-link; roles derived via `PersonService`.
- [ ] **10.9** Update person components to display profileImage and roles derived from `PersonService`.
- [ ] **10.10** Rewrite `PropertyModule`: ownership derived via `PropertyService` from new ownership collections.
- [ ] **10.11** Update property components to display owners (persons and cobrands), agents, tenants, and household members.
- [ ] **10.12** Write E2E tests for person create/edit with profileImage, role assignment, and property ownership display.

## Phase 11 — Cleanup

- [ ] **11.1** Remove frontend device-related modules, components, routes (`/user/devices/`), repositories, and services.
- [ ] **11.2** Remove `DeviceContext` from the user layout.
- [ ] **11.3** Run full test suite on the OpenBSD server (`npm run test:unit`, `npm run test:integration`, `go test ./...`) and confirm all tests pass.
- [ ] **11.4** Deploy backend to OpenBSD server and deploy frontend independently (see `CLAUDE.local.md` for deployment steps).
