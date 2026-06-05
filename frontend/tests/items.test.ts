import { test, expect } from '@playwright/test';
import PocketBase from 'pocketbase';

const PB_URL = 'http://192.168.8.144:8090';
const ADMIN_EMAIL = 'casperchen91@hotmail.com';
const ADMIN_PASSWORD = 'w3m#@tpth100';
const TEST_EMAIL = 'e2e_items@keybook.test';
const TEST_PASSWORD = 'E2Eitems_test1';

test.beforeAll(async () => {
  const pb = new PocketBase(PB_URL);
  await pb.collection('_superusers').authWithPassword(ADMIN_EMAIL, ADMIN_PASSWORD);
  try {
    await pb.collection('users').create({
      email: TEST_EMAIL,
      password: TEST_PASSWORD,
      passwordConfirm: TEST_PASSWORD,
    });
  } catch {
    // user already exists
  }
});

test.beforeEach(async ({ page }) => {
  await page.goto('/auth/login');
  await page.getByLabel('Email').fill(TEST_EMAIL);
  await page.getByLabel('Password').fill(TEST_PASSWORD);
  await page.getByText('Log in!').click();
  await page.waitForURL('/user/properties/list');
});

test('create item', async ({ page }) => {
  await page.goto('/user/items');
  await page.getByText('Add Item').click();
  await page.waitForURL('/user/items/add');
  await page.getByLabel('Item Name').fill('Front Door Key');
  await page.getByLabel('Description').fill('Key to the front door');
  await page.getByText('Save').click();
  await page.waitForURL('/user/items');
  await expect(page.getByText('Front Door Key')).toBeVisible();
});

test('edit item name', async ({ page }) => {
  const pb = new PocketBase(PB_URL);
  await pb.collection('_superusers').authWithPassword(ADMIN_EMAIL, ADMIN_PASSWORD);
  const users = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  const testUser = users[0];

  const persons = await pb.collection('persons').getFullList({ filter: `user = "${testUser.id}"` });
  let personId: string;
  if (persons.length === 0) {
    const person = await pb.collection('persons').create({ name: 'Test Person', DOB: '1990-01-01', user: testUser.id });
    personId = person.id;
  } else {
    personId = persons[0].id;
  }

  const item = await pb.collection('items').create({ name: 'Edit Me', description: 'Original description' });
  await pb.collection('personItems').create({ person: personId, item: item.id });

  await page.goto('/user/items');
  await expect(page.getByText('Edit Me')).toBeVisible();
  await page.getByText('Edit').first().click();
  await page.waitForURL('/user/items/edit');

  const nameInput = page.getByLabel('Item Name');
  await nameInput.clear();
  await nameInput.fill('Edited Key');
  await page.getByText('Save').click();
  await page.waitForURL('/user/items');
  await expect(page.getByText('Edited Key')).toBeVisible();

  await pb.collection('items').delete(item.id);
});

test('delete item', async ({ page }) => {
  const pb = new PocketBase(PB_URL);
  await pb.collection('_superusers').authWithPassword(ADMIN_EMAIL, ADMIN_PASSWORD);
  const users = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  const testUser = users[0];

  const persons = await pb.collection('persons').getFullList({ filter: `user = "${testUser.id}"` });
  let personId: string;
  if (persons.length === 0) {
    const person = await pb.collection('persons').create({ name: 'Test Person', DOB: '1990-01-01', user: testUser.id });
    personId = person.id;
  } else {
    personId = persons[0].id;
  }

  const item = await pb.collection('items').create({ name: 'Delete Me', description: 'To be deleted' });
  await pb.collection('personItems').create({ person: personId, item: item.id });

  await page.goto('/user/items');
  await expect(page.getByText('Delete Me')).toBeVisible();
  await page.getByText('Edit').first().click();
  await page.waitForURL('/user/items/edit');

  await page.getByText('Delete Item').click();
  await page.getByText('Confirm').click();
  await page.waitForURL('/user/items');
  await expect(page.getByText('Delete Me')).not.toBeVisible();
});

test('designate entry device', async ({ page }) => {
  const pb = new PocketBase(PB_URL);
  await pb.collection('_superusers').authWithPassword(ADMIN_EMAIL, ADMIN_PASSWORD);
  const users = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  const testUser = users[0];

  const persons = await pb.collection('persons').getFullList({ filter: `user = "${testUser.id}"` });
  let personId: string;
  if (persons.length === 0) {
    const person = await pb.collection('persons').create({ name: 'Test Person', DOB: '1990-01-01', user: testUser.id });
    personId = person.id;
  } else {
    personId = persons[0].id;
  }

  const item = await pb.collection('items').create({ name: 'Key With Device', description: '' });
  await pb.collection('personItems').create({ person: personId, item: item.id });

  await page.goto('/user/items');
  await page.getByText('Edit').first().click();
  await page.waitForURL('/user/items/edit');

  await page.getByLabel('Device Type').fill('Key');
  await page.getByLabel('Identifier').fill('K-001');
  await page.getByLabel('Defunct Reason').fill('None');
  await page.getByText('Create Entry Device').click();
  await page.waitForURL('/user/items');

  const entryDevices = await pb.collection('entryDevices').getFullList({ filter: `item = "${item.id}"` });
  expect(entryDevices.length).toBe(1);
  expect(entryDevices[0].identifier).toBe('K-001');

  await pb.collection('entryDevices').delete(entryDevices[0].id);
  await pb.collection('items').delete(item.id);
});

test('mark entry device defunct', async ({ page }) => {
  const pb = new PocketBase(PB_URL);
  await pb.collection('_superusers').authWithPassword(ADMIN_EMAIL, ADMIN_PASSWORD);
  const users = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  const testUser = users[0];

  const persons = await pb.collection('persons').getFullList({ filter: `user = "${testUser.id}"` });
  let personId: string;
  if (persons.length === 0) {
    const person = await pb.collection('persons').create({ name: 'Test Person', DOB: '1990-01-01', user: testUser.id });
    personId = person.id;
  } else {
    personId = persons[0].id;
  }

  const item = await pb.collection('items').create({ name: 'Defunct Key', description: '' });
  await pb.collection('personItems').create({ person: personId, item: item.id });
  const ed = await pb.collection('entryDevices').create({
    item: item.id,
    deviceType: 'Key',
    identifier: 'K-002',
    defunctReason: 'None',
  });

  await page.goto('/user/items');
  await page.getByText('Edit').first().click();
  await page.waitForURL('/user/items/edit');

  const defunctInput = page.getByLabel('Defunct Reason');
  await defunctInput.clear();
  await defunctInput.fill('Lost');
  await page.getByText('Save Entry Device').click();
  await page.waitForURL('/user/items');

  const updated = await pb.collection('entryDevices').getOne(ed.id);
  expect(updated.defunctReason).toBe('Lost');

  await pb.collection('entryDevices').delete(ed.id);
  await pb.collection('items').delete(item.id);
});
