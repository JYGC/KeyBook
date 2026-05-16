# Requirements — full-collections

## 1. Item Registry

### 1.1 Items

WHEN a user creates an item THE SYSTEM SHALL store its name, description, and optional picture.
WHEN a user edits an item THE SYSTEM SHALL update the stored name, description, and picture.
WHEN a user deletes an item THE SYSTEM SHALL remove all associated personItems, propertyItems, and entryDevices records.
WHEN a user views an item THE SYSTEM SHALL display its name, description, picture, and all associations (person owners, property locations, entry device status).

### 1.2 Person item ownership

WHEN a user assigns an item to a person THE SYSTEM SHALL create a personItems record linking that person and item.
WHEN a user assigns an item to a person it is already assigned to THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes an item from a person THE SYSTEM SHALL delete the corresponding personItems record.
WHEN a user views a person's items THE SYSTEM SHALL list all items linked to that person via personItems.

### 1.3 Property item inventory

WHEN a user assigns an item to a property THE SYSTEM SHALL create a propertyItems record linking that property and item.
WHEN a user assigns an item to a property it is already assigned to THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes an item from a property THE SYSTEM SHALL delete the corresponding propertyItems record.
WHEN a user views a property's inventory THE SYSTEM SHALL list all items linked to that property via propertyItems.

### 1.4 Entry devices

WHEN a user designates an item as an entry device THE SYSTEM SHALL create an entryDevices record with the item, a deviceType, and an identifier.
WHEN a user designates an item that already has an entryDevices record THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user changes any field of an entry device THE SYSTEM SHALL update the entryDevices record.
WHEN a user marks an entry device as defunct THE SYSTEM SHALL store a defunctReason (Lost/Damaged/Retired/Stolen) on the entryDevices record.
WHEN a user attempts to mark an entry device as defunct without providing a defunctReason THE SYSTEM SHALL display a validation error.
WHEN a user removes the entry device designation from an item THE SYSTEM SHALL delete the entryDevices record (the underlying item is retained).
WHEN a user views entry devices for a property THE SYSTEM SHALL list all items that have an entryDevices record and are assigned to that property via propertyItems.

## 2. Real Estate Agent Management

### 2.1 Agents

WHEN a user registers a person as an agent THE SYSTEM SHALL create an agents record linking that person to a cobrand.
WHEN a user registers a person as an agent for a cobrand they are already an agent of THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes an agent THE SYSTEM SHALL delete the agents record and all associated propertyAgents records.
WHEN a user views an agent THE SYSTEM SHALL display the linked person details, cobrand, and all properties the agent manages.

### 2.2 Property agents

WHEN a user assigns an agent to a property THE SYSTEM SHALL create a propertyAgents record linking that agent and property.
WHEN a user assigns an agent to a property they are already assigned to THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes an agent from a property THE SYSTEM SHALL delete the corresponding propertyAgents record.
WHEN a user views a property's agents THE SYSTEM SHALL list all agents linked to that property via propertyAgents.

## 3. Property Ownership

### 3.1 Property owners

WHEN a user creates a property owner record THE SYSTEM SHALL link it to a property.
WHEN a user links a person to a property owner THE SYSTEM SHALL create a personPropertyOwners record with the person and propertyOwner.
WHEN a user links a cobrand (firm) to a property owner THE SYSTEM SHALL create a cobrandPropertyOwners record with the cobrand and propertyOwner.
WHEN a user views ownership of a property THE SYSTEM SHALL display all persons and cobrands linked as owners via their respective relation collections.
WHEN a user removes a person from a property owner THE SYSTEM SHALL delete the corresponding personPropertyOwners record.
WHEN a user removes a cobrand from a property owner THE SYSTEM SHALL delete the corresponding cobrandPropertyOwners record.

## 4. Cobrand (Firm) Management

### 4.1 Cobrands

WHEN a user creates a cobrand THE SYSTEM SHALL store it as a named firm entity.
WHEN a user views a cobrand THE SYSTEM SHALL display its admins, agents, managed properties, and owned properties.

### 4.2 Cobrand admins

WHEN a user is added as a cobrand admin THE SYSTEM SHALL create a cobrandAdmins record linking the user to the cobrand.
WHEN a user is added as a cobrand admin for a cobrand they already administer THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a cobrand admin is removed THE SYSTEM SHALL delete the cobrandAdmins record.

### 4.3 Cobrand property managers

WHEN a cobrand is assigned to manage a property THE SYSTEM SHALL create a cobrandPropertyManagers record linking the cobrand and property.
WHEN a cobrand is assigned to manage a property it already manages THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a cobrand is removed from managing a property THE SYSTEM SHALL delete the corresponding cobrandPropertyManagers record.
WHEN a user views a cobrand's managed properties THE SYSTEM SHALL list all properties linked via cobrandPropertyManagers.

## 5. Tenant and Household Management

### 5.1 Tenants

WHEN a user adds a person as a tenant of a property THE SYSTEM SHALL create a tenants record linking that person and property.
WHEN a user adds a person as a tenant of a property they are already a tenant of THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes a tenant from a property THE SYSTEM SHALL delete the corresponding tenants record.
WHEN a user views tenants of a property THE SYSTEM SHALL list all persons linked via tenants.

### 5.2 Households

WHEN a user adds a person to a property's household THE SYSTEM SHALL create a households record linking that person and property.
WHEN a user adds a person to a household they are already a member of THE SYSTEM SHALL display an error and not create a duplicate record.
WHEN a user removes a person from a household THE SYSTEM SHALL delete the corresponding households record.
WHEN a user views a property's household THE SYSTEM SHALL list all persons linked via households.

## 6. Data model restructuring

### 6.1 Users

WHEN a user account is created THE SYSTEM SHALL store only an email address. The default `name` and `avatar` fields provided by PocketBase shall be removed.

### 6.2 Persons

WHEN a person is created THE SYSTEM SHALL store name, date of birth, an optional profile image, and an optional link to a user account.
WHEN a user links a person to a user account THE SYSTEM SHALL set the persons.user field to the given users record.
WHEN a user views a person THE SYSTEM SHALL derive the person's property roles (owner, tenant, agent, household member) from the relevant relation collections.

### 6.3 Properties

WHEN a property is created THE SYSTEM SHALL store an address.
WHEN access to a property is checked THE SYSTEM SHALL determine authorisation via propertyOwners → personPropertyOwners or cobrandPropertyOwners.
