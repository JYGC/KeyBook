import { test, expect } from '@playwright/test';
import PocketBase from 'pocketbase';

const PB_URL = 'http://192.168.8.144:8090';
const ADMIN_EMAIL = 'casperchen91@hotmail.com';
const ADMIN_PASSWORD = 'w3m#@tpth100';
const TEST_EMAIL = 'e2e_properties@keybook.test';
const TEST_PASSWORD = 'E2Eproperties_test1';

async function adminAuth(): Promise<PocketBase> {
  const backendClient = new PocketBase(PB_URL);
  const adminAuthResponse = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ identity: ADMIN_EMAIL, password: ADMIN_PASSWORD }),
  });
  const adminAuthResult = await adminAuthResponse.json();
  backendClient.authStore.save(adminAuthResult.token, adminAuthResult.admin);
  return backendClient;
}

async function ensureTestPerson(backendClient: PocketBase, userId: string): Promise<string> {
  const persons = await backendClient.collection('persons').getFullList({ filter: `user = "${userId}"` });
  if (persons.length > 0) return persons[0].id;
  const person = await backendClient.collection('persons').create<{ id: string }>({
    name: 'E2E Properties Test Person',
    DOB: '1990-01-01',
    user: userId,
  });
  return person.id;
}

async function cleanupTestTenants(backendClient: PocketBase, propertyId: string): Promise<void> {
  const tenants = await backendClient
    .collection('tenants')
    .getFullList({ filter: `property = "${propertyId}"` });
  for (const tenant of tenants) await backendClient.collection('tenants').delete(tenant.id);
}

async function cleanupTestHousehold(backendClient: PocketBase, propertyId: string): Promise<void> {
  const members = await backendClient
    .collection('households')
    .getFullList({ filter: `property = "${propertyId}"` });
  for (const member of members) await backendClient.collection('households').delete(member.id);
}

test.beforeAll(async () => {
  const backendClient = await adminAuth();
  try {
    await backendClient.collection('users').create({
      email: TEST_EMAIL,
      password: TEST_PASSWORD,
      passwordConfirm: TEST_PASSWORD,
    });
  } catch {
    // user already exists
  }
  const users = await backendClient.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  if (users.length > 0) await ensureTestPerson(backendClient, users[0].id);
});

test.beforeEach(async ({ page }) => {
  await page.goto('/auth/login');
  await page.getByLabel('Email').fill(TEST_EMAIL);
  await page.getByLabel('Password').fill(TEST_PASSWORD);
  await page.getByText('Log in!').click();
  await page.waitForURL(/\/user\/properties\/list/);
});

test('create property', async ({ page }) => {
  const backendClient = await adminAuth();
  const uniqueAddress = `1 E2E Properties Test St ${Date.now()}`;

  await page.goto('/user/properties/list');
  await page.waitForLoadState('networkidle');
  await page.getByRole('button', { name: 'Add Property' }).click();
  await page.waitForURL(/\/user\/properties\/add/);
  await page.waitForLoadState('networkidle');

  await page.getByLabel('Property Address').fill(uniqueAddress);
  await page.getByRole('button', { name: 'Save' }).click();
  await page.waitForURL(/\/user\/properties\/list/);

  const properties = await backendClient.collection('properties').getFullList({
    filter: `address = "${uniqueAddress}"`,
  });
  expect(properties.length).toBe(1);
  // Skip property/propertyOwner deletion — last-owner hook prevents it
});

test('add and remove tenant', async ({ page }) => {
  const backendClient = await adminAuth();

  const users = await backendClient.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  const testUserId = users[0].id;
  const ownerPersonId = await ensureTestPerson(backendClient, testUserId);

  const property = await backendClient.collection('properties').create<{ id: string }>({ address: '2 E2E Properties Tenant Test Ave' });
  const propertyOwner = await backendClient.collection('propertyOwners').create<{ id: string }>({ property: property.id });
  await backendClient.collection('personPropertyOwners').create({ propertyOwner: propertyOwner.id, person: ownerPersonId });

  const tenantPerson = await backendClient.collection('persons').create<{ id: string }>({
    name: 'E2E Properties Tenant Person',
    DOB: '1995-07-10',
  });

  await page.goto(`/user/properties/detail?id=${property.id}`);
  await page.waitForLoadState('networkidle');

  await page.getByLabel('Person ID').first().fill(tenantPerson.id);
  await page.getByRole('button', { name: 'Add Tenant' }).click();
  await page.waitForTimeout(1000);

  const tenants = await backendClient
    .collection('tenants')
    .getFullList({ filter: `property = "${property.id}"` });
  expect(tenants.some((tenant) => tenant.person === tenantPerson.id)).toBe(true);

  await page.getByRole('button', { name: 'Remove' }).first().click();
  await page.waitForTimeout(1000);

  const tenantsAfter = await backendClient
    .collection('tenants')
    .getFullList({ filter: `property = "${property.id}"` });
  expect(tenantsAfter.length).toBe(0);

  // Skip property/propertyOwner — the last-owner hook guard blocks deletion.
  await backendClient.collection('persons').delete(tenantPerson.id);
  const personPropertyOwners = await backendClient.collection('personPropertyOwners').getFullList({ filter: `propertyOwner = "${propertyOwner.id}"` });
  for (const personPropertyOwner of personPropertyOwners) await backendClient.collection('personPropertyOwners').delete(personPropertyOwner.id);
});

test('add household member', async ({ page }) => {
  const backendClient = await adminAuth();

  const users = await backendClient.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  const testUserId = users[0].id;
  const ownerPersonId = await ensureTestPerson(backendClient, testUserId);

  const property = await backendClient.collection('properties').create<{ id: string }>({ address: '3 E2E Properties Household Test Ave' });
  const propertyOwner = await backendClient.collection('propertyOwners').create<{ id: string }>({ property: property.id });
  await backendClient.collection('personPropertyOwners').create({ propertyOwner: propertyOwner.id, person: ownerPersonId });

  const householdPerson = await backendClient.collection('persons').create<{ id: string }>({
    name: 'E2E Properties Household Person',
    DOB: '2000-12-01',
  });

  await page.goto(`/user/properties/detail?id=${property.id}`);
  await page.waitForLoadState('networkidle');

  // Find the household member Person ID input (second TextInput with label "Person ID")
  await page.getByLabel('Person ID').nth(1).fill(householdPerson.id);
  await page.getByRole('button', { name: 'Add Household Member' }).click();
  await page.waitForTimeout(1000);

  const members = await backendClient
    .collection('households')
    .getFullList({ filter: `property = "${property.id}"` });
  expect(members.some((member) => member.person === householdPerson.id)).toBe(true);

  await cleanupTestHousehold(backendClient, property.id);
  await backendClient.collection('persons').delete(householdPerson.id);
  const personPropertyOwners = await backendClient.collection('personPropertyOwners').getFullList({ filter: `propertyOwner = "${propertyOwner.id}"` });
  for (const personPropertyOwner of personPropertyOwners) await backendClient.collection('personPropertyOwners').delete(personPropertyOwner.id);
});
