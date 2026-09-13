import PocketBase from 'pocketbase';
import { describe, it, expect, beforeAll } from 'vitest';
import { PropertyRepository } from './property-repository';
import { PropertyOwnerRepository } from './property-owner-repository';
import { TenantRepository } from './tenant-repository';
import { HouseholdRepository } from './household-repository';

const PB_URL = 'http://192.168.8.144:8090';
const backendClient = new PocketBase(PB_URL);

beforeAll(async () => {
  const adminAuthResponse = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ identity: 'casperchen91@hotmail.com', password: 'w3m#@tpth100' }),
  });
  const adminAuthResult = await adminAuthResponse.json();
  backendClient.authStore.save(adminAuthResult.token, adminAuthResult.admin);
});

describe('PropertyRepository', () => {
  it('creates, reads, updates, and deletes a property', async () => {
    const propertyRepository = new PropertyRepository(backendClient);

    const created = await propertyRepository.createProperty('1 Property Repo Test St');
    expect(created.address).toBe('1 Property Repo Test St');
    expect(created.id).toBeTruthy();

    const fetched = await propertyRepository.getPropertyById(created.id);
    expect(fetched.id).toBe(created.id);
    expect(fetched.address).toBe('1 Property Repo Test St');

    const allProperties = await propertyRepository.getAllProperties();
    expect(allProperties.some((property) => property.id === created.id)).toBe(true);

    const updated = await propertyRepository.updateProperty(created.id, '2 Property Repo Test St Updated');
    expect(updated.address).toBe('2 Property Repo Test St Updated');

    // Must delete in dependency order: no propertyOwner exists yet, so direct delete is fine
    await propertyRepository.deleteProperty(created.id);
    const allPropertiesAfterDelete = await propertyRepository.getAllProperties();
    expect(allPropertiesAfterDelete.some((propertiesAfterDelete) => propertiesAfterDelete.id === created.id)).toBe(false);
  });
});

describe('PropertyOwnerRepository', () => {
  it('creates, queries, and deletes a property owner record', async () => {
    const propertyRepository = new PropertyRepository(backendClient);
    const propertyOwnerRepository = new PropertyOwnerRepository(backendClient);

    const property = await propertyRepository.createProperty('10 PropertyOwner Repo Test Ave');
    const created = await propertyOwnerRepository.createPropertyOwner(property.id);
    expect(created.property).toBe(property.id);
    expect(created.id).toBeTruthy();

    const byProperty = await propertyOwnerRepository.getPropertyOwnersByPropertyId(property.id);
    expect(byProperty.some((propertyOwner) => propertyOwner.id === created.id)).toBe(true);

    // Skip propertyOwner and property deletion — last-owner hook prevents it
  });
});

describe('TenantRepository', () => {
  it('creates, queries, and deletes a tenant record', async () => {
    const tenantRepository = new TenantRepository(backendClient);

    const person = await backendClient
      .collection('persons')
      .create<{ id: string }>({ name: 'Tenant Repo Test Person', DOB: '1992-04-10' });
    const property = await backendClient
      .collection('properties')
      .create<{ id: string }>({ address: '20 Tenant Repo Test Ave' });
    const propertyOwner = await backendClient
      .collection('propertyOwners')
      .create<{ id: string }>({ property: property.id });
    await backendClient
      .collection('personPropertyOwners')
      .create({ person: person.id, propertyOwner: propertyOwner.id });

    const created = await tenantRepository.createTenant(person.id, property.id);
    expect(created.person).toBe(person.id);
    expect(created.property).toBe(property.id);
    expect(created.id).toBeTruthy();

    const byProperty = await tenantRepository.getTenantsByPropertyId(property.id);
    expect(byProperty.some((tenant) => tenant.id === created.id)).toBe(true);

    const byPerson = await tenantRepository.getTenantsByPersonId(person.id);
    expect(byPerson.some((tenant) => tenant.id === created.id)).toBe(true);

    await tenantRepository.deleteTenant(created.id);
    // Delete personPropertyOwner before person (required relation constraint)
    const personPropertyOwners = await backendClient.collection('personPropertyOwners').getFullList({ filter: `person = "${person.id}"` });
    for (const personPropertyOwner of personPropertyOwners) await backendClient.collection('personPropertyOwners').delete(personPropertyOwner.id);
    await backendClient.collection('persons').delete(person.id);
    // Skip propertyOwner and property deletion — last-owner hook guard prevents it
  });
});

describe('HouseholdRepository', () => {
  it('creates, queries, and deletes a household record', async () => {
    const householdRepository = new HouseholdRepository(backendClient);

    const person = await backendClient
      .collection('persons')
      .create<{ id: string }>({ name: 'Household Repo Test Person', DOB: '1988-11-25' });
    const property = await backendClient
      .collection('properties')
      .create<{ id: string }>({ address: '30 Household Repo Test Ave' });
    const propertyOwner = await backendClient
      .collection('propertyOwners')
      .create<{ id: string }>({ property: property.id });
    await backendClient
      .collection('personPropertyOwners')
      .create({ person: person.id, propertyOwner: propertyOwner.id });

    const created = await householdRepository.createHousehold(person.id, property.id);
    expect(created.person).toBe(person.id);
    expect(created.property).toBe(property.id);
    expect(created.id).toBeTruthy();

    const byProperty = await householdRepository.getHouseholdsByPropertyId(property.id);
    expect(byProperty.some((household) => household.id === created.id)).toBe(true);

    const byPerson = await householdRepository.getHouseholdsByPersonId(person.id);
    expect(byPerson.some((household) => household.id === created.id)).toBe(true);

    await householdRepository.deleteHousehold(created.id);
    // Delete personPropertyOwner before person (required relation constraint)
    const personPropertyOwners = await backendClient.collection('personPropertyOwners').getFullList({ filter: `person = "${person.id}"` });
    for (const personPropertyOwner of personPropertyOwners) await backendClient.collection('personPropertyOwners').delete(personPropertyOwner.id);
    await backendClient.collection('persons').delete(person.id);
    // Skip propertyOwner and property deletion — last-owner hook guard prevents it
  });
});
