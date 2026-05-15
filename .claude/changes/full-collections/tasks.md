# Tasks — full-collections

## Phase 1 — Schema

- [ ] **1.1** Write unit tests for schema validation helpers that verify required fields and relation targets for all collections.
- [ ] **1.2** Create a PocketBase migration in `backend/migrations/` that defines all collections as specified in `design.md` section 2.1, including carry-over columns (`users.avatar`, `persons.profileImage`, `entryDevices.identifier`, `entryDevices.defunctReason`). Delete `database/pb_schema.json`.
- [ ] **1.3** In the same migration, set PocketBase access rules for all collections using the new ownership chain (`propertyOwners → personPropertyOwners → person.user` and `cobrandPropertyOwners → cobrand → cobrandAdmins → userId`).

## Phase 2 — Backend: DTOs

- [ ] **2.1** Write unit tests for all DTOs covering field mapping and required-field validation.
- [ ] **2.2** Add DTOs: `CobrandDto`, `CobrandAdminDto`, `CobrandPropertyManagerDto`, `CobrandPropertyOwnerDto`.
- [ ] **2.3** Add DTOs: `PropertyOwnerDto`, `PersonPropertyOwnerDto`.
- [ ] **2.4** Add DTOs: `AgentDto`, `PropertyAgentDto`.
- [ ] **2.5** Add DTOs: `HouseholdDto`, `TenantDto`.
- [ ] **2.6** Add DTOs: `ItemDto`, `PropertyItemDto`, `PersonItemDto`, `EntryDeviceDto` (include `identifier` and `defunctReason`).
- [ ] **2.7** Rewrite `PersonDto`: fields are `name`, `DOB`, `user`, `profileImage`.
- [ ] **2.8** Rewrite `PropertyDto`: fields are `address` only.

## Phase 3 — Backend: Repositories

- [ ] **3.1** Write unit tests for each repository covering CRUD operations and relation queries.
- [ ] **3.2** Add `ItemRepository`, `PropertyItemRepository`, `PersonItemRepository`, `EntryDeviceRepository`.
- [ ] **3.3** Add `PropertyOwnerRepository`, `PersonPropertyOwnerRepository`.
- [ ] **3.4** Add `CobrandRepository`, `CobrandAdminRepository`, `CobrandPropertyManagerRepository`, `CobrandPropertyOwnerRepository`.
- [ ] **3.5** Add `AgentRepository`, `PropertyAgentRepository`.
- [ ] **3.6** Add `HouseholdRepository`, `TenantRepository`.
- [ ] **3.7** Rewrite `PersonRepository`: fields are `name`, `DOB`, `user`, `profileImage`.
- [ ] **3.8** Rewrite `PropertyRepository`: fields are `address` only.
- [ ] **3.9** Delete `DeviceRepository` and `PersonDeviceRepository`.

## Phase 4 — Backend: DI wiring

- [ ] **4.1** Register all new repositories in the `dig` container in `cmd/keybook.go`.
- [ ] **4.2** Remove `DeviceHistoryServices` and `PersonDeviceHistoryServices` hooks and their registrations from `cmd/keybook.go`.
- [ ] **4.3** Remove `DeviceRepository` and `PersonDeviceRepository` from the DI container.

## Phase 5 — Frontend: Item registry

- [ ] **5.1** Write unit tests for `ItemListModule` and `ItemDetailModule` covering CRUD, entry device toggling, identifier, and defunctReason.
- [ ] **5.2** Implement `src/lib/modules/item/` — `ItemListModule`, `ItemDetailModule`.
- [ ] **5.3** Add `ItemContext` to the user layout.
- [ ] **5.4** Implement components: `ItemList`, `ItemEditor`, `EntryDeviceEditor` (includes identifier, deviceType, defunctReason).
- [ ] **5.5** Add routes: `/user/items/`, `/user/items/add/`, `/user/items/edit/`.
- [ ] **5.6** Write E2E tests for item create, edit, delete, entry device designation, and defunct marking flows.

## Phase 6 — Frontend: Cobrand management

- [ ] **6.1** Write unit tests for `CobrandListModule` and `CobrandDetailModule`.
- [ ] **6.2** Implement `src/lib/modules/cobrand/` — `CobrandListModule`, `CobrandDetailModule`.
- [ ] **6.3** Add `CobrandContext` to the user layout.
- [ ] **6.4** Implement components: `CobrandList`, `CobrandEditor`, `CobrandAdminList`, `CobrandPropertyManagerList`.
- [ ] **6.5** Add routes: `/user/cobrands/`, `/user/cobrands/add/`, `/user/cobrands/[id]/`.
- [ ] **6.6** Write E2E tests for cobrand create, admin assignment, and property manager assignment.

## Phase 7 — Frontend: Agent management

- [ ] **7.1** Write unit tests for `AgentListModule` and `AgentDetailModule`.
- [ ] **7.2** Implement `src/lib/modules/agent/` — `AgentListModule`, `AgentDetailModule`.
- [ ] **7.3** Add `AgentContext` to the user layout.
- [ ] **7.4** Implement components: `AgentList`, `AgentEditor`, `PropertyAgentList`.
- [ ] **7.5** Add routes: `/user/agents/`, `/user/agents/add/`, `/user/agents/[id]/`.
- [ ] **7.6** Write E2E tests for agent registration and property assignment.

## Phase 8 — Frontend: Persons and properties

- [ ] **8.1** Write unit tests for updated `PersonModule` and `PropertyModule`.
- [ ] **8.2** Rewrite `PersonModule`: fields are name, DOB, profileImage, user-link; roles (owner/tenant/agent/household) derived from relation collections.
- [ ] **8.3** Update person components to display profileImage and roles derived from relation collections.
- [ ] **8.4** Rewrite `PropertyModule`: ownership derived from `propertyOwners`/`personPropertyOwners`/`cobrandPropertyOwners`; no direct owners/managers fields.
- [ ] **8.5** Update property components to display owners (persons and cobrands), agents, tenants, and household members.
- [ ] **8.6** Write E2E tests for person create/edit with profileImage, role assignment, and property ownership display.

## Phase 9 — Cleanup

- [ ] **9.1** Remove frontend device-related modules, components, and routes (`/user/devices/`).
- [ ] **9.2** Remove `DeviceContext` from the user layout.
- [ ] **9.3** Run full test suite (`npm run test`, `go test ./...`) and confirm all tests pass.
- [ ] **9.4** Deploy to OpenBSD server and verify the running application (see CLAUDE.md deployment steps).
