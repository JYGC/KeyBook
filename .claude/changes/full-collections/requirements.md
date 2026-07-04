# Requirements — full-collections

## 1. Item Registry

### 1.1 Items

WHEN a user creates an item THE SYSTEM SHALL store its name, description, and optional picture.
WHEN a user edits an item THE SYSTEM SHALL update the stored name, description, and picture.
WHEN a user deletes an item THE SYSTEM SHALL remove all associated personItems, propertyItems, and entryDevices records.
WHEN a user views an item THE SYSTEM SHALL display its name, description, picture, and all associations (person owners, property locations, entry device status).
WHEN an unauthenticated request attempts to list, view, create, update, or delete an items record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view an items record THE SYSTEM SHALL allow it only if they are a person-owner of that item, a person owner or cobrand owner or cobrand manager of a property it is assigned to, or a tenant, household member, or agent at such a property.
WHEN an authenticated user attempts to create an items record THE SYSTEM SHALL allow it regardless of existing associations.
WHEN an authenticated user attempts to update or delete an items record THE SYSTEM SHALL allow it only if they are a person-owner of that item, or a person owner, cobrand owner, or cobrand manager of a property it is assigned to.

### 1.2 Person item ownership

WHEN a user assigns an item to a person THE SYSTEM SHALL create a personItems record linking that person and item.
WHEN a user assigns an item to a person it is already assigned to THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes an item from a person THE SYSTEM SHALL delete the corresponding personItems record.
WHEN a user views a person's items THE SYSTEM SHALL list all items linked to that person via personItems.
WHEN an unauthenticated request attempts to list, view, create, update, or delete a personItems record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a personItems record THE SYSTEM SHALL allow it only if they are the linked person, or a person owner, cobrand owner, or cobrand manager of a property where that item is located.
WHEN an authenticated user attempts to create a personItems record THE SYSTEM SHALL allow it regardless of existing associations.
WHEN an authenticated user attempts to update or delete a personItems record THE SYSTEM SHALL allow it only if they are the linked person, or a person owner, cobrand owner, or cobrand manager of a property where that item is located.

### 1.3 Property item inventory

WHEN a user assigns an item to a property THE SYSTEM SHALL create a propertyItems record linking that property and item.
WHEN a user assigns an item to a property it is already assigned to THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes an item from a property THE SYSTEM SHALL delete the corresponding propertyItems record.
WHEN a user views a property's inventory THE SYSTEM SHALL list all items linked to that property via propertyItems.
WHEN an unauthenticated request attempts to list, view, create, update, or delete a propertyItems record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a propertyItems record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, cobrand manager, tenant, household member, or agent of the linked property, or a person-owner of the linked item.
WHEN an authenticated user attempts to create, update, or delete a propertyItems record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, or cobrand manager of the linked property, or a person-owner of the linked item.

### 1.4 Entry devices

WHEN a user designates an item as an entry device THE SYSTEM SHALL create an entryDevices record with the item, a deviceType, and an identifier.
WHEN a user designates an item that already has an entryDevices record THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user changes any field of an entry device THE SYSTEM SHALL update the entryDevices record.
WHEN a user marks an entry device as defunct THE SYSTEM SHALL store a defunctReason (Lost/Damaged/Retired/Stolen) on the entryDevices record.
WHEN a user attempts to mark an entry device as defunct without providing a defunctReason THE SYSTEM SHALL display a validation error.
WHEN a user removes the entry device designation from an item THE SYSTEM SHALL delete the entryDevices record (the underlying item is retained).
WHEN a user views entry devices for a property THE SYSTEM SHALL list all items that have an entryDevices record and are assigned to that property via propertyItems.
WHEN an unauthenticated request attempts to list, view, create, update, or delete an entryDevices record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view an entryDevices record THE SYSTEM SHALL allow it only if they are a person-owner of the linked item, or a person owner, cobrand owner, cobrand manager, tenant, household member, or agent of a property the item is assigned to.
WHEN an authenticated user attempts to create, update, or delete an entryDevices record THE SYSTEM SHALL allow it only if they are a person-owner of the linked item, or a person owner, cobrand owner, or cobrand manager of a property the item is assigned to.

## 2. Real Estate Agent Management

### 2.1 Agents

WHEN a user registers a person as an agent THE SYSTEM SHALL create an agents record linking that person to a cobrand.
WHEN a user registers a person as an agent for a cobrand they are already an agent of THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes an agent THE SYSTEM SHALL delete the agents record and all associated propertyAgents records.
WHEN a user views an agent THE SYSTEM SHALL display the linked person details, cobrand, and all properties the agent manages.
WHEN an unauthenticated request attempts to list, view, create, update, or delete an agents record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view an agents record THE SYSTEM SHALL allow it only if they are a cobrand admin of that agent's cobrand, the agent themselves, or an owner or manager of a property that agent is assigned to.
WHEN an authenticated user attempts to create or update an agents record THE SYSTEM SHALL allow it only if they are a cobrand admin of the linked cobrand.
WHEN an authenticated user attempts to delete an agents record THE SYSTEM SHALL allow it if they are a cobrand admin of the linked cobrand, or if the record links their own person account.

### 2.2 Property agents

WHEN a user assigns an agent to a property THE SYSTEM SHALL create a propertyAgents record linking that agent and property.
WHEN a user assigns an agent to a property they are already assigned to THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes an agent from a property THE SYSTEM SHALL delete the corresponding propertyAgents record.
WHEN a user views a property's agents THE SYSTEM SHALL list all agents linked to that property via propertyAgents.
WHEN an unauthenticated request attempts to list, view, create, update, or delete a propertyAgents record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a propertyAgents record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, or cobrand manager of the linked property, or if they are the assigned agent or a cobrand admin of the agent's cobrand.
WHEN an authenticated user attempts to create or update a propertyAgents record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, or cobrand manager of the linked property, or a cobrand admin of the agent's cobrand.
WHEN an authenticated user attempts to delete a propertyAgents record THE SYSTEM SHALL allow it if they are a person owner, cobrand owner, cobrand manager, or cobrand admin of the agent's cobrand, or if the record links their own person account as the agent.

## 3. Property Ownership

### 3.1 Property owners

WHEN a user creates a property owner record THE SYSTEM SHALL link it to a property.
WHEN a user links a person to a property owner THE SYSTEM SHALL create a personPropertyOwners record with the person and propertyOwner.
WHEN a user links a cobrand (firm) to a property owner THE SYSTEM SHALL create a cobrandPropertyOwners record with the cobrand and propertyOwner.
WHEN a user views ownership of a property THE SYSTEM SHALL display all persons and cobrands linked as owners via their respective relation collections.
WHEN a user removes a person from a property owner THE SYSTEM SHALL delete the corresponding personPropertyOwners record.
WHEN a user removes a cobrand from a property owner THE SYSTEM SHALL delete the corresponding cobrandPropertyOwners record.
WHEN an unauthenticated request attempts to list, view, create, update, or delete a propertyOwners record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a propertyOwners record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, or cobrand manager of the linked property.
WHEN a tenant, household member, or agent attempts to list or view a propertyOwners record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to create a propertyOwners record for a property that already has at least one owner and they are not themselves an existing owner of that property THE SYSTEM SHALL reject the request.
WHEN an authenticated user attempts to update or delete a propertyOwners record THE SYSTEM SHALL allow it only if they are a person owner or cobrand owner of the linked property.
WHEN a cobrand manager attempts to update or delete a propertyOwners record THE SYSTEM SHALL reject it with 403.
WHEN a user attempts to delete the last remaining propertyOwners record for a property THE SYSTEM SHALL reject the request.

### 3.2 Person property owners

WHEN an unauthenticated request attempts to list, view, create, update, or delete a personPropertyOwners record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a personPropertyOwners record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, or cobrand manager of the linked property.
WHEN a tenant, household member, or agent attempts to list or view a personPropertyOwners record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to create, update, or delete a personPropertyOwners record THE SYSTEM SHALL allow it only if they are a person owner or cobrand owner of the linked property.
WHEN a cobrand manager attempts to create, update, or delete a personPropertyOwners record THE SYSTEM SHALL reject it with 403.

### 3.3 Cobrand property owners

WHEN an unauthenticated request attempts to list, view, create, update, or delete a cobrandPropertyOwners record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a cobrandPropertyOwners record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, or cobrand manager of the linked property.
WHEN a tenant, household member, or agent attempts to list or view a cobrandPropertyOwners record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to create or update a cobrandPropertyOwners record THE SYSTEM SHALL allow it only if they are a person owner or cobrand owner of the linked property.
WHEN a cobrand manager attempts to create or update a cobrandPropertyOwners record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to delete a cobrandPropertyOwners record THE SYSTEM SHALL allow it if they are a person owner or cobrand owner of the linked property, or if they are an admin of the cobrand being removed.

## 4. Cobrand (Firm) Management

### 4.1 Cobrands

WHEN a user creates a cobrand THE SYSTEM SHALL store it as a named firm entity.
WHEN a user views a cobrand THE SYSTEM SHALL display its admins, agents, managed properties, and owned properties.
WHEN an unauthenticated request attempts to list, view, create, update, or delete a cobrands record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a cobrands record THE SYSTEM SHALL allow it only if they are an admin of that cobrand or an agent belonging to that cobrand.
WHEN an authenticated user who is neither a cobrand admin nor an agent of that cobrand attempts to list or view a cobrands record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to create a cobrands record THE SYSTEM SHALL allow it regardless of existing cobrand membership.
WHEN an authenticated user attempts to update or delete a cobrands record THE SYSTEM SHALL allow it only if they are an admin of that cobrand.
WHEN a cobrand agent attempts to update or delete a cobrands record THE SYSTEM SHALL reject it with 403.

### 4.2 Cobrand admins

WHEN a user is added as a cobrand admin THE SYSTEM SHALL create a cobrandAdmins record linking the user to the cobrand.
WHEN a user is added as a cobrand admin for a cobrand they already administer THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a cobrand admin is removed THE SYSTEM SHALL delete the cobrandAdmins record.
WHEN an unauthenticated request attempts to list, view, create, update, or delete a cobrandAdmins record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a cobrandAdmins record THE SYSTEM SHALL allow it only if they are an existing admin of that cobrand or an agent belonging to that cobrand.
WHEN an authenticated user attempts to create a cobrandAdmins record THE SYSTEM SHALL allow it regardless of existing admin membership.
WHEN an authenticated user attempts to create a cobrandAdmins record for a cobrand that already has at least one admin and they are not themselves an existing admin of that cobrand THE SYSTEM SHALL reject the request.
WHEN an authenticated user attempts to update a cobrandAdmins record THE SYSTEM SHALL allow it only if they are an existing admin of that cobrand.
WHEN an authenticated user attempts to delete a cobrandAdmins record THE SYSTEM SHALL allow it if they are an existing admin of that cobrand, or if the record links their own user account.
WHEN a cobrand agent attempts to update or delete a cobrandAdmins record THE SYSTEM SHALL reject it with 403.

### 4.3 Cobrand property managers

WHEN a cobrand is assigned to manage a property THE SYSTEM SHALL create a cobrandPropertyManagers record linking the cobrand and property.
WHEN a cobrand is assigned to manage a property it already manages THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a cobrand is removed from managing a property THE SYSTEM SHALL delete the corresponding cobrandPropertyManagers record.
WHEN a user views a cobrand's managed properties THE SYSTEM SHALL list all properties linked via cobrandPropertyManagers.
WHEN an unauthenticated request attempts to list, view, create, update, or delete a cobrandPropertyManagers record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a cobrandPropertyManagers record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, or an admin of the managing cobrand.
WHEN an authenticated user attempts to create or update a cobrandPropertyManagers record THE SYSTEM SHALL allow it only if they are a person owner or cobrand owner of the linked property.
WHEN a managing cobrand admin attempts to create or update a cobrandPropertyManagers record for a property they do not own THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to delete a cobrandPropertyManagers record THE SYSTEM SHALL allow it if they are a person owner or cobrand owner of the linked property, or if they are an admin of the managing cobrand.

## 5. Tenant and Household Management

### 5.1 Tenants

WHEN a user adds a person as a tenant of a property THE SYSTEM SHALL create a tenants record linking that person and property.
WHEN a user adds a person as a tenant of a property they are already a tenant of THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes a tenant from a property THE SYSTEM SHALL delete the corresponding tenants record.
WHEN a user views tenants of a property THE SYSTEM SHALL list all persons linked via tenants.
WHEN an unauthenticated request attempts to list, view, create, update, or delete a tenants record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a tenants record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, or cobrand manager of the linked property, or a tenant or household member of that property.
WHEN an authenticated user attempts to create or update a tenants record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, or cobrand manager of the linked property.
WHEN an authenticated user attempts to delete a tenants record THE SYSTEM SHALL allow it if they are a person owner, cobrand owner, or cobrand manager of the linked property, or if the record links their own person account.

### 5.2 Households

WHEN a user adds a person to a property's household THE SYSTEM SHALL create a households record linking that person and property.
WHEN a user adds a person to a household they are already a member of THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes a person from a household THE SYSTEM SHALL delete the corresponding households record.
WHEN a user views a property's household THE SYSTEM SHALL list all persons linked via households.
WHEN an unauthenticated request attempts to list, view, create, update, or delete a households record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a households record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, or cobrand manager of the linked property, or a household member or tenant of that property.
WHEN an authenticated user attempts to create or update a households record THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, or cobrand manager of the linked property.
WHEN an authenticated user attempts to delete a households record THE SYSTEM SHALL allow it if they are a person owner, cobrand owner, or cobrand manager of the linked property, or if the record links their own person account.

## 6. Data model restructuring

### 6.1 Users

WHEN a user account is created THE SYSTEM SHALL store only an email address. The default `name` and `avatar` fields provided by PocketBase shall be removed.

### 6.2 Persons

WHEN a person is created THE SYSTEM SHALL store name, date of birth, an optional profile image, and an optional link to a user account.
WHEN a user links a person to a user account THE SYSTEM SHALL set the persons.user field to the given users record.
WHEN a user views a person THE SYSTEM SHALL derive the person's property roles (owner, tenant, agent, household member) from the relevant relation collections.
WHEN an unauthenticated request attempts to list, view, or create a person record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a person record THE SYSTEM SHALL allow it only if at least one of the following is true: the person's user field matches the authenticated user; the authenticated user owns or manages a property where that person is a tenant; the authenticated user owns or manages a property that holds an item currently assigned to that person via personItems.
WHEN an authenticated user attempts to list or view a person record that does not satisfy any of those conditions THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to create a person record THE SYSTEM SHALL allow it regardless of ownership relation, so that owners can register tenants who do not yet have user accounts.
WHEN an authenticated user attempts to update a person record THE SYSTEM SHALL allow it only if that person's user field is linked to the authenticated user's account.
WHEN an authenticated user attempts to delete a person record THE SYSTEM SHALL allow it only if that person's user field is linked to the authenticated user's account.
WHEN an authenticated user attempts to update or delete a person record whose user field is null or linked to a different user THE SYSTEM SHALL reject it with 403.

### 6.3 Properties

WHEN a property is created THE SYSTEM SHALL store an address.
WHEN access to a property is checked THE SYSTEM SHALL determine authorisation via propertyOwners → personPropertyOwners or cobrandPropertyOwners.
WHEN an unauthenticated request attempts to list, view, create, update, or delete a property record THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to list or view a property THE SYSTEM SHALL allow it only if they are a person owner, cobrand owner, cobrand manager, tenant, household member, or agent of that property.
WHEN an authenticated user attempts to list or view a property and none of those relations hold THE SYSTEM SHALL reject it with 403.
WHEN an authenticated user attempts to create a property record THE SYSTEM SHALL allow it regardless of existing ownership.
WHEN an authenticated user attempts to update a property record THE SYSTEM SHALL allow it only if they are a person owner or cobrand owner of that property.
WHEN an authenticated user attempts to delete a property record THE SYSTEM SHALL allow it only if they are a person owner or cobrand owner of that property.
WHEN a cobrand manager, tenant, household member, or agent attempts to update or delete a property record THE SYSTEM SHALL reject it with 403.

## 7. User onboarding

### 7.1 Registration creates a person

WHEN a user submits the registration form THE SYSTEM SHALL collect name, date of birth, email, password, and password confirmation.
WHEN a user submits a valid registration form THE SYSTEM SHALL create the user account (email only), authenticate the new user, and create a persons record with the given name and date of birth whose user field links to the new account.
WHEN the user account cannot be created (e.g. email already registered) THE SYSTEM SHALL display the error and not create a persons record.
WHEN the persons record cannot be created after the user account has been created THE SYSTEM SHALL still complete login and rely on post-login routing to offer person setup.
WHEN registration completes THE SYSTEM SHALL navigate to /user and apply post-login routing.

### 7.2 Post-login routing

WHEN an authenticated user arrives at /user and a persons record linked to their account exists THE SYSTEM SHALL redirect to the property list page.
WHEN an authenticated user arrives at /user and no persons record is linked to their account THE SYSTEM SHALL redirect to the person setup page.
WHEN the linked-person lookup fails THE SYSTEM SHALL display an error instead of redirecting.

### 7.3 Person setup page

WHEN a user with no linked person opens the person setup page THE SYSTEM SHALL display a form collecting name and date of birth.
WHEN a user submits a valid person setup form THE SYSTEM SHALL create a persons record linked to their account and navigate to the property list page.
WHEN a user submits a person setup form with a missing name or date of birth THE SYSTEM SHALL display a validation error and not create a record.
WHEN a user who already has a linked person opens the person setup page THE SYSTEM SHALL redirect to the property list page.

