import PocketBase from 'pocketbase';
import { describe, it, expect, beforeAll } from 'vitest';
import { ItemRepository } from './item-repository';
import { EntryDeviceRepository } from './entry-device-repository';
import { PropertyItemRepository } from './property-item-repository';
import { PersonItemRepository } from './person-item-repository';

const backendClient = new PocketBase('http://192.168.8.144:8090');

beforeAll(async () => {
  // PocketBase Go v0.22 uses /api/admins (not _superusers which is v0.23+).
  // PocketBase JS SDK v0.26 maps backendClient.admins to _superusers, so we use raw fetch.
  const adminAuthResponse = await fetch('http://192.168.8.144:8090/api/admins/auth-with-password', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ identity: 'casperchen91@hotmail.com', password: 'w3m#@tpth100' }),
  });
  const adminAuthResult = await adminAuthResponse.json();
  backendClient.authStore.save(adminAuthResult.token, adminAuthResult.admin);
});

describe('ItemRepository', () => {
  it('creates, reads, updates, and deletes an item', async () => {
    const itemRepository = new ItemRepository(backendClient);

    const created = await itemRepository.createItem('Test Key', 'A test key item');
    expect(created.name).toBe('Test Key');
    expect(created.description).toBe('A test key item');
    expect(created.id).toBeTruthy();

    const fetched = await itemRepository.getItemById(created.id);
    expect(fetched.id).toBe(created.id);
    expect(fetched.name).toBe('Test Key');

    await itemRepository.updateItem(created.id, 'Updated Key', 'Updated description');
    const updated = await itemRepository.getItemById(created.id);
    expect(updated.name).toBe('Updated Key');
    expect(updated.description).toBe('Updated description');

    const allItems = await itemRepository.getAllItems();
    expect(allItems.some((item) => item.id === created.id)).toBe(true);

    await itemRepository.deleteItem(created.id);
  });
});

describe('EntryDeviceRepository', () => {
  it('creates, reads, updates, and deletes an entry device', async () => {
    const itemRepository = new ItemRepository(backendClient);
    const entryDeviceRepository = new EntryDeviceRepository(backendClient);

    const item = await itemRepository.createItem('Entry Device Test Item', 'desc');

    const created = await entryDeviceRepository.createEntryDevice(item.id, 'Key', 'K-001', 'None');
    expect(created.item).toBe(item.id);
    expect(created.deviceType).toBe('Key');
    expect(created.identifier).toBe('K-001');
    expect(created.defunctReason).toBe('None');

    const byItem = await entryDeviceRepository.getEntryDeviceByItemId(item.id);
    expect(byItem?.id).toBe(created.id);

    const entryDeviceFetchedById = await entryDeviceRepository.getEntryDeviceById(created.id);
    expect(entryDeviceFetchedById.id).toBe(created.id);

    await entryDeviceRepository.updateEntryDevice(created.id, 'Fob', 'F-002', 'None');
    const updated = await entryDeviceRepository.getEntryDeviceByItemId(item.id);
    expect(updated?.deviceType).toBe('Fob');
    expect(updated?.identifier).toBe('F-002');

    await entryDeviceRepository.deleteEntryDevice(created.id);
    const afterDelete = await entryDeviceRepository.getEntryDeviceByItemId(item.id);
    expect(afterDelete).toBeNull();

    await itemRepository.deleteItem(item.id);
  });
});

describe('PropertyItemRepository', () => {
  it('creates, queries, and deletes a property-item link', async () => {
    const itemRepository = new ItemRepository(backendClient);
    const propertyItemRepository = new PropertyItemRepository(backendClient);

    const item = await itemRepository.createItem('Property Item Test', 'desc');
    const property = await backendClient.collection('properties').create<{ id: string }>({ address: '99 Test Ave' });

    const created = await propertyItemRepository.createPropertyItem(item.id, property.id);
    expect(created.item).toBe(item.id);
    expect(created.property).toBe(property.id);

    const byItem = await propertyItemRepository.getPropertyItemsByItemId(item.id);
    expect(byItem.some((propertyItem) => propertyItem.id === created.id)).toBe(true);

    const byProperty = await propertyItemRepository.getPropertyItemsByPropertyId(property.id);
    expect(byProperty.some((propertyItem) => propertyItem.id === created.id)).toBe(true);

    await propertyItemRepository.deletePropertyItem(created.id);
    await itemRepository.deleteItem(item.id);
    await backendClient.collection('properties').delete(property.id);
  });
});

describe('PersonItemRepository', () => {
  it('creates, queries, and deletes a person-item link', async () => {
    const itemRepository = new ItemRepository(backendClient);
    const personItemRepository = new PersonItemRepository(backendClient);

    const item = await itemRepository.createItem('Person Item Test', 'desc');
    const person = await backendClient
      .collection('persons')
      .create<{ id: string }>({ name: 'Test Person', DOB: '1990-01-01 00:00:00.000Z' });

    const created = await personItemRepository.createPersonItem(person.id, item.id);
    expect(created.item).toBe(item.id);
    expect(created.person).toBe(person.id);

    const byItem = await personItemRepository.getPersonItemsByItemId(item.id);
    expect(byItem.some((personItem) => personItem.id === created.id)).toBe(true);

    const byPerson = await personItemRepository.getPersonItemsByPersonId(person.id);
    expect(byPerson.some((personItem) => personItem.id === created.id)).toBe(true);

    await personItemRepository.deletePersonItem(created.id);
    await itemRepository.deleteItem(item.id);
    await backendClient.collection('persons').delete(person.id);
  });
});
