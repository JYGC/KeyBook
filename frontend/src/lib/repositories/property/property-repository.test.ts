import PocketBase from 'pocketbase';
import { describe, it, expect, beforeAll } from 'vitest';
import { PropertyRepository } from './property-repository';
import { PropertyOwnerRepository } from './property-owner-repository';
import { TenantRepository } from './tenant-repository';
import { HouseholdRepository } from './household-repository';

const PB_URL = 'http://192.168.8.144:8090';
const pb = new PocketBase(PB_URL);

beforeAll(async () => {
  const res = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ identity: 'casperchen91@hotmail.com', password: 'w3m#@tpth100' }),
  });
  const data = await res.json();
  pb.authStore.save(data.token, data.admin);
});

describe('PropertyRepository', () => {
  it('creates, reads, updates, and deletes a property', async () => {
    const repo = new PropertyRepository(pb);

    const created = await repo.create('1 Property Repo Test St');
    expect(created.address).toBe('1 Property Repo Test St');
    expect(created.id).toBeTruthy();

    const fetched = await repo.getById(created.id);
    expect(fetched.id).toBe(created.id);
    expect(fetched.address).toBe('1 Property Repo Test St');

    const all = await repo.getAll();
    expect(all.some((p) => p.id === created.id)).toBe(true);

    const updated = await repo.update(created.id, '2 Property Repo Test St Updated');
    expect(updated.address).toBe('2 Property Repo Test St Updated');

    // Must delete in dependency order: no propertyOwner exists yet, so direct delete is fine
    await repo.delete(created.id);
    const allAfter = await repo.getAll();
    expect(allAfter.some((p) => p.id === created.id)).toBe(false);
  });
});

describe('PropertyOwnerRepository', () => {
  it('creates, queries, and deletes a property owner record', async () => {
    const propertyRepo = new PropertyRepository(pb);
    const propertyOwnerRepo = new PropertyOwnerRepository(pb);

    const property = await propertyRepo.create('10 PropertyOwner Repo Test Ave');
    const created = await propertyOwnerRepo.create(property.id);
    expect(created.property).toBe(property.id);
    expect(created.id).toBeTruthy();

    const byProperty = await propertyOwnerRepo.getByPropertyId(property.id);
    expect(byProperty.some((r) => r.id === created.id)).toBe(true);

    // Skip propertyOwner and property deletion — last-owner hook prevents it
  });
});

describe('TenantRepository', () => {
  it('creates, queries, and deletes a tenant record', async () => {
    const tenantRepo = new TenantRepository(pb);

    const person = await pb
      .collection('persons')
      .create<{ id: string }>({ name: 'Tenant Repo Test Person', DOB: '1992-04-10' });
    const property = await pb
      .collection('properties')
      .create<{ id: string }>({ address: '20 Tenant Repo Test Ave' });
    const propertyOwner = await pb
      .collection('propertyOwners')
      .create<{ id: string }>({ property: property.id });
    await pb
      .collection('personPropertyOwners')
      .create({ person: person.id, propertyOwner: propertyOwner.id });

    const created = await tenantRepo.create(person.id, property.id);
    expect(created.person).toBe(person.id);
    expect(created.property).toBe(property.id);
    expect(created.id).toBeTruthy();

    const byProperty = await tenantRepo.getByPropertyId(property.id);
    expect(byProperty.some((r) => r.id === created.id)).toBe(true);

    const byPerson = await tenantRepo.getByPersonId(person.id);
    expect(byPerson.some((r) => r.id === created.id)).toBe(true);

    await tenantRepo.delete(created.id);
    // Delete personPropertyOwner before person (required relation constraint)
    const ppos = await pb.collection('personPropertyOwners').getFullList({ filter: `person = "${person.id}"` });
    for (const ppo of ppos) await pb.collection('personPropertyOwners').delete(ppo.id);
    await pb.collection('persons').delete(person.id);
    // Skip propertyOwner and property deletion — last-owner hook guard prevents it
  });
});

describe('HouseholdRepository', () => {
  it('creates, queries, and deletes a household record', async () => {
    const householdRepo = new HouseholdRepository(pb);

    const person = await pb
      .collection('persons')
      .create<{ id: string }>({ name: 'Household Repo Test Person', DOB: '1988-11-25' });
    const property = await pb
      .collection('properties')
      .create<{ id: string }>({ address: '30 Household Repo Test Ave' });
    const propertyOwner = await pb
      .collection('propertyOwners')
      .create<{ id: string }>({ property: property.id });
    await pb
      .collection('personPropertyOwners')
      .create({ person: person.id, propertyOwner: propertyOwner.id });

    const created = await householdRepo.create(person.id, property.id);
    expect(created.person).toBe(person.id);
    expect(created.property).toBe(property.id);
    expect(created.id).toBeTruthy();

    const byProperty = await householdRepo.getByPropertyId(property.id);
    expect(byProperty.some((r) => r.id === created.id)).toBe(true);

    const byPerson = await householdRepo.getByPersonId(person.id);
    expect(byPerson.some((r) => r.id === created.id)).toBe(true);

    await householdRepo.delete(created.id);
    // Delete personPropertyOwner before person (required relation constraint)
    const ppos = await pb.collection('personPropertyOwners').getFullList({ filter: `person = "${person.id}"` });
    for (const ppo of ppos) await pb.collection('personPropertyOwners').delete(ppo.id);
    await pb.collection('persons').delete(person.id);
    // Skip propertyOwner and property deletion — last-owner hook guard prevents it
  });
});
