import { test, expect } from '@playwright/test';
import PocketBase from 'pocketbase';

const PB_URL = 'http://192.168.8.144:8090';
const ADMIN_EMAIL = 'casperchen91@hotmail.com';
const ADMIN_PASSWORD = 'w3m#@tpth100';
const TEST_EMAIL = 'e2e_cobrands@keybook.test';
const TEST_PASSWORD = 'E2Ecobrands_test1';
const TEST_EMAIL_2 = 'e2e_cobrands2@keybook.test';
const TEST_PASSWORD_2 = 'E2Ecobrands2_test1';

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

async function ensureTestPerson(pb: PocketBase, userId: string): Promise<string> {
  const persons = await pb.collection('persons').getFullList({ filter: `user = "${userId}"` });
  if (persons.length > 0) return persons[0].id;
  const person = await pb.collection('persons').create({
    name: 'E2E Cobrands Test Person',
    DOB: '1990-01-01',
    user: userId,
  });
  return person.id;
}

async function cleanupTestCobrands(pb: PocketBase): Promise<void> {
  const testNames = ['E2E Test Cobrand', 'E2E Admin Test Cobrand', 'E2E Manager Test Cobrand'];
  for (const name of testNames) {
    const found = await pb.collection('cobrands').getFullList({ filter: `name = "${name}"` });
    for (const c of found) {
      const admins = await pb.collection('cobrandAdmins').getFullList({ filter: `cobrand = "${c.id}"` });
      for (const a of admins) await pb.collection('cobrandAdmins').delete(a.id);
      const managers = await pb.collection('cobrandPropertyManagers').getFullList({ filter: `cobrand = "${c.id}"` });
      for (const m of managers) await pb.collection('cobrandPropertyManagers').delete(m.id);
      await pb.collection('cobrands').delete(c.id);
    }
  }
}

test.beforeAll(async () => {
  const pb = await adminAuth();
  for (const [email, password] of [
    [TEST_EMAIL, TEST_PASSWORD],
    [TEST_EMAIL_2, TEST_PASSWORD_2],
  ]) {
    try {
      await pb.collection('users').create({
        email,
        password,
        passwordConfirm: password,
      });
    } catch {
      // user already exists
    }
  }
  await cleanupTestCobrands(pb);
});

test.beforeEach(async ({ page }) => {
  await page.goto('/auth/login');
  await page.getByLabel('Email').fill(TEST_EMAIL);
  await page.getByLabel('Password').fill(TEST_PASSWORD);
  await page.getByText('Log in!').click();
  await page.waitForURL(/\/user\/properties\/list/);
});

test('create cobrand', async ({ page }) => {
  const pb = await adminAuth();

  await page.goto('/user/cobrands');
  await page.waitForLoadState('networkidle');
  await page.getByRole('button', { name: 'Add Cobrand' }).click();
  await page.waitForURL(/\/user\/cobrands\/add/);
  await page.getByLabel('Cobrand Name').fill('E2E Test Cobrand');
  await page.getByRole('button', { name: 'Save' }).click();
  await page.waitForURL(/\/user\/cobrands(\/)?$/);

  const cobrands = await pb
    .collection('cobrands')
    .getFullList({ filter: 'name = "E2E Test Cobrand"' });
  expect(cobrands.length).toBe(1);

  await cleanupTestCobrands(pb);
});

test('add admin to cobrand', async ({ page }) => {
  const pb = await adminAuth();

  const users1 = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  const users2 = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL_2}"` });
  const testUserId = users1[0].id;
  const secondUserId = users2[0].id;

  const cobrand = await pb.collection('cobrands').create<{ id: string }>({ name: 'E2E Admin Test Cobrand' });
  await pb.collection('cobrandAdmins').create({ user: testUserId, cobrand: cobrand.id });

  await page.goto(`/user/cobrands/detail?id=${cobrand.id}`);
  await page.waitForLoadState('networkidle');

  await page.getByLabel('User ID').fill(secondUserId);
  await page.getByRole('button', { name: 'Add Admin' }).click();
  await page.waitForTimeout(1000);

  const admins = await pb
    .collection('cobrandAdmins')
    .getFullList({ filter: `cobrand = "${cobrand.id}"` });
  expect(admins.some((a) => a.user === secondUserId)).toBe(true);

  await cleanupTestCobrands(pb);
});

test('add property manager to cobrand', async ({ page }) => {
  const pb = await adminAuth();

  const users1 = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  const testUserId = users1[0].id;
  const personId = await ensureTestPerson(pb, testUserId);

  const cobrand = await pb.collection('cobrands').create<{ id: string }>({ name: 'E2E Manager Test Cobrand' });
  await pb.collection('cobrandAdmins').create({ user: testUserId, cobrand: cobrand.id });

  const property = await pb.collection('properties').create<{ id: string }>({ address: '77 E2E Manager Ave' });
  const propertyOwner = await pb.collection('propertyOwners').create<{ id: string }>({ property: property.id });
  await pb.collection('personPropertyOwners').create({ propertyOwner: propertyOwner.id, person: personId });

  await page.goto(`/user/cobrands/detail?id=${cobrand.id}`);
  await page.waitForLoadState('networkidle');

  await page.getByLabel('Property ID').fill(property.id);
  await page.getByRole('button', { name: 'Add Property Manager' }).click();
  await page.waitForTimeout(1000);

  const managers = await pb
    .collection('cobrandPropertyManagers')
    .getFullList({ filter: `cobrand = "${cobrand.id}"` });
  expect(managers.some((m) => m.property === property.id)).toBe(true);

  // Cleanup cobrandPropertyManagers and cobrand
  for (const m of managers) await pb.collection('cobrandPropertyManagers').delete(m.id);
  await cleanupTestCobrands(pb);
  // Clean up personPropertyOwner; skip propertyOwner+property due to last-owner hook guard
  const ppos = await pb.collection('personPropertyOwners').getFullList({ filter: `propertyOwner = "${propertyOwner.id}"` });
  for (const ppo of ppos) await pb.collection('personPropertyOwners').delete(ppo.id);
});
