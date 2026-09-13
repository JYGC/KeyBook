import { test, expect } from '@playwright/test';
import PocketBase from 'pocketbase';

const PB_URL = 'http://192.168.8.144:8090';
const ADMIN_EMAIL = 'casperchen91@hotmail.com';
const ADMIN_PASSWORD = 'w3m#@tpth100';
const TEST_EMAIL = 'e2e_agents@keybook.test';
const TEST_PASSWORD = 'E2Eagents_test1';

async function adminAuth(): Promise<PocketBase> {
	const backendClient = new PocketBase(PB_URL);
	const adminAuthResponse = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ identity: ADMIN_EMAIL, password: ADMIN_PASSWORD })
	});
	const adminAuthResult = await adminAuthResponse.json();
	backendClient.authStore.save(adminAuthResult.token, adminAuthResult.admin);
	return backendClient;
}

async function ensureTestPerson(backendClient: PocketBase, userId: string): Promise<string> {
	const persons = await backendClient.collection('persons').getFullList({ filter: `user = "${userId}"` });
	if (persons.length > 0) return persons[0].id;
	const person = await backendClient.collection('persons').create({
		name: 'E2E Agents Test Person',
		DOB: '1990-01-01',
		user: userId
	});
	return person.id;
}

async function ensureTestCobrand(backendClient: PocketBase, userId: string): Promise<string> {
	const found = await backendClient
		.collection('cobrands')
		.getFullList({ filter: 'name = "E2E Agents Test Cobrand"' });
	if (found.length > 0) return found[0].id;
	const cobrand = await backendClient
		.collection('cobrands')
		.create<{ id: string }>({ name: 'E2E Agents Test Cobrand' });
	// approved: true — this fixture stands in for the cobrand admin performing agent
	// registration, which Phase 12's approval gate requires the admin to be approved for.
	await backendClient
		.collection('cobrandAdmins')
		.create({ user: userId, cobrand: cobrand.id, approved: true });
	return cobrand.id;
}

async function cleanupTestAgents(backendClient: PocketBase): Promise<void> {
	const agents = await backendClient.collection('agents').getFullList();
	for (const agent of agents) {
		const propertyAgents = await backendClient.collection('propertyAgents').getFullList({ filter: `agent = "${a.id}"` });
		for (const propertyAgent of propertyAgents) await backendClient.collection('propertyAgents').delete(propertyAgent.id);
		await backendClient.collection('agents').delete(a.id);
	}
	const found = await backendClient
		.collection('cobrands')
		.getFullList({ filter: 'name = "E2E Agents Test Cobrand"' });
	for (const cobrand of found) {
		const admins = await backendClient
			.collection('cobrandAdmins')
			.getFullList({ filter: `cobrand = "${c.id}"` });
		for (const admin of admins) await backendClient.collection('cobrandAdmins').delete(admin.id);
		await backendClient.collection('cobrands').delete(c.id);
	}
}

test.beforeAll(async () => {
	const backendClient = await adminAuth();
	try {
		await backendClient.collection('users').create({
			email: TEST_EMAIL,
			password: TEST_PASSWORD,
			passwordConfirm: TEST_PASSWORD
		});
	} catch {
		// user already exists
	}
	await cleanupTestAgents(backendClient);
});

test.beforeEach(async ({ page }) => {
	await page.goto('/auth/login');
	await page.getByLabel('Email').fill(TEST_EMAIL);
	await page.getByLabel('Password').fill(TEST_PASSWORD);
	await page.getByText('Log in!').click();
	// Post-login routing depends on whether this user has a linked person/cobrandAdmins
	// record at the time of login (see full-collections Phase 13); these tests don't
	// depend on which of the three destinations it lands on.
	await page.waitForURL(/\/user\/(properties\/list|cobrands|setup)/);
});

test('register agent', async ({ page }) => {
	const backendClient = await adminAuth();

	const users = await backendClient.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
	const testUserId = users[0].id;
	const personId = await ensureTestPerson(backendClient, testUserId);
	const cobrandId = await ensureTestCobrand(backendClient, testUserId);

	await page.goto('/user/agents');
	await page.waitForLoadState('networkidle');
	await page.getByRole('button', { name: 'Add Agent' }).click();
	await page.waitForURL(/\/user\/agents\/add/);

	await page.getByLabel('Person ID').fill(personId);
	await page.getByLabel('Cobrand ID').fill(cobrandId);
	await page.getByRole('button', { name: 'Save' }).click();
	await page.waitForURL(/\/user\/agents(\/)?$/);

	const agents = await backendClient
		.collection('agents')
		.getFullList({ filter: `person = "${personId}" && cobrand = "${cobrandId}"` });
	expect(agents.length).toBe(1);

	await cleanupTestAgents(backendClient);
});

test('add property assignment to agent', async ({ page }) => {
	const backendClient = await adminAuth();

	const users = await backendClient.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
	const testUserId = users[0].id;
	const personId = await ensureTestPerson(backendClient, testUserId);
	const cobrandId = await ensureTestCobrand(backendClient, testUserId);

	const agent = await backendClient
		.collection('agents')
		.create<{ id: string }>({ person: personId, cobrand: cobrandId });
	const property = await backendClient
		.collection('properties')
		.create<{ id: string }>({ address: '88 E2E Agent Ave' });
	const propertyOwner = await backendClient
		.collection('propertyOwners')
		.create<{ id: string }>({ property: property.id });
	await backendClient
		.collection('personPropertyOwners')
		.create({ propertyOwner: propertyOwner.id, person: personId });

	await page.goto(`/user/agents/detail?id=${agent.id}`);
	await page.waitForLoadState('networkidle');

	await page.getByLabel('Property ID').fill(property.id);
	await page.getByRole('button', { name: 'Add Property Assignment' }).click();
	await page.waitForTimeout(1000);

	const propertyAgents = await backendClient
		.collection('propertyAgents')
		.getFullList({ filter: `agent = "${agent.id}"` });
	expect(propertyAgents.some((propertyAgent) => propertyAgent.property === property.id)).toBe(true);

	// Skip propertyOwner/property — the last-owner hook guard blocks deletion.
	for (const propertyAgent of propertyAgents) await backendClient.collection('propertyAgents').delete(propertyAgent.id);
	await cleanupTestAgents(backendClient);
	const personPropertyOwners = await backendClient
		.collection('personPropertyOwners')
		.getFullList({ filter: `propertyOwner = "${propertyOwner.id}"` });
	for (const personPropertyOwner of personPropertyOwners) await backendClient.collection('personPropertyOwners').delete(personPropertyOwner.id);
});
