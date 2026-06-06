import { test, expect } from '@playwright/test';
import PocketBase from 'pocketbase';

const PB_URL = 'http://192.168.8.144:8090';
const ADMIN_EMAIL = 'casperchen91@hotmail.com';
const ADMIN_PASSWORD = 'w3m#@tpth100';
const TEST_EMAIL = 'e2e_items@keybook.test';
const TEST_PASSWORD = 'E2Eitems_test1';

async function adminAuth(): Promise<PocketBase> {
  const pb = new PocketBase(PB_URL);
  const res = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ identity: ADMIN_EMAIL, password: ADMIN_PASSWORD }),
  });
  const data = await res.json();
  pb.authStore.save(data.token, data.admin);
  return pb;
}

async function ensureTestPerson(pb: PocketBase): Promise<string> {
  const users = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  const testUser = users[0];
  const persons = await pb.collection('persons').getFullList({ filter: `user = "${testUser.id}"` });
  if (persons.length > 0) return persons[0].id;
  const person = await pb.collection('persons').create({
    name: 'E2E Test Person',
    DOB: '1990-01-01',
    user: testUser.id,
  });
  return person.id;
}

async function cleanupTestUserItems(pb: PocketBase): Promise<void> {
  const users = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  if (users.length === 0) return;
  const persons = await pb.collection('persons').getFullList({ filter: `user = "${users[0].id}"` });
  for (const person of persons) {
    const pis = await pb.collection('personItems').getFullList({ filter: `person = "${person.id}"` });
    for (const pi of pis) {
      const itemId = pi.item as string;
      try {
        const eds = await pb.collection('entryDevices').getFullList({ filter: `item = "${itemId}"` });
        for (const ed of eds) await pb.collection('entryDevices').delete(ed.id);
      } catch { /* ignore */ }
      try { await pb.collection('personItems').delete(pi.id); } catch { /* ignore */ }
      try { await pb.collection('items').delete(itemId); } catch { /* ignore */ }
    }
  }
}

test.beforeAll(async () => {
  const pb = await adminAuth();
  try {
    await pb.collection('users').create({
      email: TEST_EMAIL,
      password: TEST_PASSWORD,
      passwordConfirm: TEST_PASSWORD,
    });
  } catch {
    // user already exists
  }
  await cleanupTestUserItems(pb);
});

test.beforeEach(async ({ page }) => {
  await page.goto('/auth/login');
  await page.getByLabel('Email').fill(TEST_EMAIL);
  await page.getByLabel('Password').fill(TEST_PASSWORD);
  await page.getByText('Log in!').click();
  await page.waitForURL(/\/user\/properties\/list/);
});

test('create item', async ({ page }) => {
  const pb = await adminAuth();

  await page.goto('/user/items');
  await page.waitForLoadState('networkidle');
  await page.getByRole('button', { name: 'Add Item' }).click();
  await page.waitForURL(/\/user\/items\/add/);
  await page.getByLabel('Item Name').fill('Front Door Key');
  await page.getByLabel('Description').fill('Key to the front door');
  await page.getByRole('button', { name: 'Save' }).click();
  await page.waitForURL(/\/user\/items(\/)?$/);

  const items = await pb.collection('items').getFullList({ filter: 'name = "Front Door Key"' });
  expect(items.length).toBe(1);
  expect(items[0].description).toBe('Key to the front door');
  await pb.collection('items').delete(items[0].id);
});

test('edit item name', async ({ page }) => {
  const pb = await adminAuth();
  const personId = await ensureTestPerson(pb);
  const item = await pb.collection('items').create({ name: 'E2E Edit Item', description: 'Original' });
  await pb.collection('personItems').create({ person: personId, item: item.id });

  await page.goto('/user/items');
  await page.waitForLoadState('networkidle');
  await expect(page.getByRole('cell', { name: 'E2E Edit Item' })).toBeVisible();

  await page.goto(`/user/items/edit?id=${item.id}`);
  await page.waitForLoadState('networkidle');

  const nameInput = page.getByLabel('Item Name');
  await nameInput.clear();
  await nameInput.fill('E2E Edited Key');
  await page.getByRole('button', { name: 'Save' }).click();
  await page.waitForURL(/\/user\/items(\/)?$/);
  await expect(page.getByRole('cell', { name: 'E2E Edited Key' })).toBeVisible();

  const items = await pb.collection('items').getFullList({ filter: 'name = "E2E Edited Key"' });
  await pb.collection('personItems').delete((await pb.collection('personItems').getFullList({ filter: `item = "${items[0].id}"` }))[0].id);
  await pb.collection('items').delete(items[0].id);
});

test('delete item', async ({ page }) => {
  const pb = await adminAuth();
  const personId = await ensureTestPerson(pb);
  const item = await pb.collection('items').create({ name: 'E2E Delete Item', description: 'Gone soon' });
  await pb.collection('personItems').create({ person: personId, item: item.id });

  page.on('dialog', async dialog => { console.log('DIALOG:', dialog.message()); await dialog.dismiss(); });

  await page.goto('/user/items');
  await page.waitForLoadState('networkidle');
  await expect(page.getByRole('cell', { name: 'E2E Delete Item' })).toBeVisible();

  await page.goto(`/user/items/edit?id=${item.id}`);
  await page.waitForLoadState('networkidle');

  await page.getByRole('button', { name: 'Delete Item' }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await page.waitForURL(/\/user\/items(\/)?$/);
  await expect(page.getByRole('cell', { name: 'E2E Delete Item' })).not.toBeVisible();
});

test('designate entry device', async ({ page }) => {
  const pb = await adminAuth();
  const personId = await ensureTestPerson(pb);
  const item = await pb
    .collection('items')
    .create({ name: 'E2E Key With Device', description: 'A test key' });
  await pb.collection('personItems').create({ person: personId, item: item.id });

  await page.goto('/user/items');
  await page.waitForLoadState('networkidle');
  await expect(page.getByRole('cell', { name: 'E2E Key With Device' })).toBeVisible();

  await page.goto(`/user/items/edit?id=${item.id}`);
  await page.waitForLoadState('networkidle');

  await page.getByLabel('Device Type').fill('Key');
  await page.getByLabel('Identifier').fill('K-001');
  await page.getByLabel('Defunct Reason').fill('None');
  await page.getByRole('button', { name: 'Create Entry Device' }).click();
  await page.waitForURL(/\/user\/items(\/)?$/);

  const entryDevices = await pb.collection('entryDevices').getFullList({ filter: `item = "${item.id}"` });
  expect(entryDevices.length).toBe(1);
  expect(entryDevices[0].identifier).toBe('K-001');

  await pb.collection('entryDevices').delete(entryDevices[0].id);
  await pb.collection('personItems').delete((await pb.collection('personItems').getFullList({ filter: `item = "${item.id}"` }))[0].id);
  await pb.collection('items').delete(item.id);
});

test('mark entry device defunct', async ({ page }) => {
  const pb = await adminAuth();
  const personId = await ensureTestPerson(pb);
  const item = await pb
    .collection('items')
    .create({ name: 'E2E Defunct Key', description: 'A key going defunct' });
  await pb.collection('personItems').create({ person: personId, item: item.id });
  const ed = await pb.collection('entryDevices').create({
    item: item.id,
    deviceType: 'Key',
    identifier: 'K-002',
    defunctReason: 'None',
  });

  await page.goto('/user/items');
  await page.waitForLoadState('networkidle');
  await expect(page.getByRole('cell', { name: 'E2E Defunct Key' })).toBeVisible();

  await page.goto(`/user/items/edit?id=${item.id}`);
  await page.waitForLoadState('networkidle');

  const defunctInput = page.getByLabel('Defunct Reason');
  await defunctInput.clear();
  await defunctInput.fill('Lost');
  await page.getByRole('button', { name: 'Save Entry Device' }).click();
  await page.waitForURL(/\/user\/items(\/)?$/);

  const updated = await pb.collection('entryDevices').getOne(ed.id);
  expect(updated.defunctReason).toBe('Lost');

  await pb.collection('entryDevices').delete(ed.id);
  await pb.collection('personItems').delete((await pb.collection('personItems').getFullList({ filter: `item = "${item.id}"` }))[0].id);
  await pb.collection('items').delete(item.id);
});
