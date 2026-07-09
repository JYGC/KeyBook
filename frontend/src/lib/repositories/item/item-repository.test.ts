import PocketBase from 'pocketbase';
import { describe, it, expect, beforeAll } from 'vitest';
import { ItemRepository } from './item-repository';
import { EntryDeviceRepository } from './entry-device-repository';
import { PropertyItemRepository } from './property-item-repository';
import { PersonItemRepository } from './person-item-repository';

const pb = new PocketBase('http://192.168.8.144:8090');

beforeAll(async () => {
  // PocketBase Go v0.22 uses /api/admins (not _superusers which is v0.23+).
  // PocketBase JS SDK v0.26 maps pb.admins to _superusers, so we use raw fetch.
  const res = await fetch('http://192.168.8.144:8090/api/admins/auth-with-password', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ identity: 'casperchen91@hotmail.com', password: 'w3m#@tpth100' }),
  });
  const data = await res.json();
  pb.authStore.save(data.token, data.admin);
});

describe('ItemRepository', () => {
  it('creates, reads, updates, and deletes an item', async () => {
    const repo = new ItemRepository(pb);

    const created = await repo.create('Test Key', 'A test key item');
    expect(created.name).toBe('Test Key');
    expect(created.description).toBe('A test key item');
    expect(created.id).toBeTruthy();

    const fetched = await repo.getById(created.id);
    expect(fetched.id).toBe(created.id);
    expect(fetched.name).toBe('Test Key');

    await repo.update(created.id, 'Updated Key', 'Updated description');
    const updated = await repo.getById(created.id);
    expect(updated.name).toBe('Updated Key');
    expect(updated.description).toBe('Updated description');

    const all = await repo.getAll();
    expect(all.some((i) => i.id === created.id)).toBe(true);

    await repo.delete(created.id);
  });
});

describe('EntryDeviceRepository', () => {
  it('creates, reads, updates, and deletes an entry device', async () => {
    const itemRepo = new ItemRepository(pb);
    const entryDeviceRepo = new EntryDeviceRepository(pb);

    const item = await itemRepo.create('Entry Device Test Item', 'desc');

    const created = await entryDeviceRepo.create(item.id, 'Key', 'K-001', 'None');
    expect(created.item).toBe(item.id);
    expect(created.deviceType).toBe('Key');
    expect(created.identifier).toBe('K-001');
    expect(created.defunctReason).toBe('None');

    const byItem = await entryDeviceRepo.getByItemId(item.id);
    expect(byItem?.id).toBe(created.id);

    const byId = await entryDeviceRepo.getById(created.id);
    expect(byId.id).toBe(created.id);

    await entryDeviceRepo.update(created.id, 'Fob', 'F-002', 'None');
    const updated = await entryDeviceRepo.getByItemId(item.id);
    expect(updated?.deviceType).toBe('Fob');
    expect(updated?.identifier).toBe('F-002');

    await entryDeviceRepo.delete(created.id);
    const afterDelete = await entryDeviceRepo.getByItemId(item.id);
    expect(afterDelete).toBeNull();

    await itemRepo.delete(item.id);
  });
});

describe('PropertyItemRepository', () => {
  it('creates, queries, and deletes a property-item link', async () => {
    const itemRepo = new ItemRepository(pb);
    const propertyItemRepo = new PropertyItemRepository(pb);

    const item = await itemRepo.create('Property Item Test', 'desc');
    const property = await pb.collection('properties').create<{ id: string }>({ address: '99 Test Ave' });

    const created = await propertyItemRepo.create(item.id, property.id);
    expect(created.item).toBe(item.id);
    expect(created.property).toBe(property.id);

    const byItem = await propertyItemRepo.getByItemId(item.id);
    expect(byItem.some((pi) => pi.id === created.id)).toBe(true);

    const byProperty = await propertyItemRepo.getByPropertyId(property.id);
    expect(byProperty.some((pi) => pi.id === created.id)).toBe(true);

    await propertyItemRepo.delete(created.id);
    await itemRepo.delete(item.id);
    await pb.collection('properties').delete(property.id);
  });
});

describe('PersonItemRepository', () => {
  it('creates, queries, and deletes a person-item link', async () => {
    const itemRepo = new ItemRepository(pb);
    const personItemRepo = new PersonItemRepository(pb);

    const item = await itemRepo.create('Person Item Test', 'desc');
    const person = await pb
      .collection('persons')
      .create<{ id: string }>({ name: 'Test Person', DOB: '1990-01-01 00:00:00.000Z' });

    const created = await personItemRepo.create(person.id, item.id);
    expect(created.item).toBe(item.id);
    expect(created.person).toBe(person.id);

    const byItem = await personItemRepo.getByItemId(item.id);
    expect(byItem.some((pi) => pi.id === created.id)).toBe(true);

    const byPerson = await personItemRepo.getByPersonId(person.id);
    expect(byPerson.some((pi) => pi.id === created.id)).toBe(true);

    await personItemRepo.delete(created.id);
    await itemRepo.delete(item.id);
    await pb.collection('persons').delete(person.id);
  });
});
