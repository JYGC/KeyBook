import { test, expect } from '@playwright/test';
import PocketBase from 'pocketbase';

const PB_URL = 'http://192.168.8.144:8090';
const ADMIN_EMAIL = 'casperchen91@hotmail.com';
const ADMIN_PASSWORD = 'w3m#@tpth100';
const TEST_EMAIL = 'e2e_agents@keybook.test';
const TEST_PASSWORD = 'E2Eagents_test1';

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
    name: 'E2E Agents Test Person',
    DOB: '1990-01-01',
    user: userId,
  });
  return person.id;
}

async function ensureTestCobrand(pb: PocketBase, userId: string): Promise<string> {
  const found = await pb.collection('cobrands').getFullList({ filter: 'name = "E2E Agents Test Cobrand"' });
  if (found.length > 0) return found[0].id;
  const cobrand = await pb.collection('cobrands').create<{ id: string }>({ name: 'E2E Agents Test Cobrand' });
  await pb.collection('cobrandAdmins').create({ user: userId, cobrand: cobrand.id });
  return cobrand.id;
}

async function cleanupTestAgents(pb: PocketBase): Promise<void> {
  const agents = await pb.collection('agents').getFullList();
  for (const a of agents) {
    const pas = await pb.collection('propertyAgents').getFullList({ filter: `agent = "${a.id}"` });
    for (const pa of pas) await pb.collection('propertyAgents').delete(pa.id);
    await pb.collection('agents').delete(a.id);
  }
  const found = await pb.collection('cobrands').getFullList({ filter: 'name = "E2E Agents Test Cobrand"' });
  for (const c of found) {
    const admins = await pb.collection('cobrandAdmins').getFullList({ filter: `cobrand = "${c.id}"` });
    for (const a of admins) await pb.collection('cobrandAdmins').delete(a.id);
    await pb.collection('cobrands').delete(c.id);
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
  await cleanupTestAgents(pb);
});

test.beforeEach(async ({ page }) => {
  await page.goto('/auth/login');
  await page.getByLabel('Email').fill(TEST_EMAIL);
  await page.getByLabel('Password').fill(TEST_PASSWORD);
  await page.getByText('Log in!').click();
  await page.waitForURL(/\/user\/properties\/list/);
});

test('register agent', async ({ page }) => {
  const pb = await adminAuth();

  const users = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  const testUserId = users[0].id;
  const personId = await ensureTestPerson(pb, testUserId);
  const cobrandId = await ensureTestCobrand(pb, testUserId);

  await page.goto('/user/agents');
  await page.waitForLoadState('networkidle');
  await page.getByRole('button', { name: 'Add Agent' }).click();
  await page.waitForURL(/\/user\/agents\/add/);

  await page.getByLabel('Person ID').fill(personId);
  await page.getByLabel('Cobrand ID').fill(cobrandId);
  await page.getByRole('button', { name: 'Save' }).click();
  await page.waitForURL(/\/user\/agents(\/)?$/);

  const agents = await pb.collection('agents').getFullList({ filter: `person = "${personId}" && cobrand = "${cobrandId}"` });
  expect(agents.length).toBe(1);

  await cleanupTestAgents(pb);
});

test('add property assignment to agent', async ({ page }) => {
  const pb = await adminAuth();

  const users = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
  const testUserId = users[0].id;
  const personId = await ensureTestPerson(pb, testUserId);
  const cobrandId = await ensureTestCobrand(pb, testUserId);

  const agent = await pb.collection('agents').create<{ id: string }>({ person: personId, cobrand: cobrandId });
  const property = await pb.collection('properties').create<{ id: string }>({ address: '88 E2E Agent Ave' });
  const propertyOwner = await pb.collection('propertyOwners').create<{ id: string }>({ property: property.id });
  await pb.collection('personPropertyOwners').create({ propertyOwner: propertyOwner.id, person: personId });

  await page.goto(`/user/agents/detail?id=${agent.id}`);
  await page.waitForLoadState('networkidle');

  await page.getByLabel('Property ID').fill(property.id);
  await page.getByRole('button', { name: 'Add Property Assignment' }).click();
  await page.waitForTimeout(1000);

  const propertyAgents = await pb
    .collection('propertyAgents')
    .getFullList({ filter: `agent = "${agent.id}"` });
  expect(propertyAgents.some((pa) => pa.property === property.id)).toBe(true);

  // Cleanup propertyAgents and agents; skip propertyOwner+property due to last-owner hook guard
  for (const pa of propertyAgents) await pb.collection('propertyAgents').delete(pa.id);
  await cleanupTestAgents(pb);
  const ppos = await pb.collection('personPropertyOwners').getFullList({ filter: `propertyOwner = "${propertyOwner.id}"` });
  for (const ppo of ppos) await pb.collection('personPropertyOwners').delete(ppo.id);
});
