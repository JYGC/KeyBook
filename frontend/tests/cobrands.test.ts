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
		name: 'E2E Cobrands Test Person',
		DOB: '1990-01-01',
		user: userId
	});
	return person.id;
}

async function cleanupTestCobrands(backendClient: PocketBase): Promise<void> {
	const testNames = ['E2E Test Cobrand', 'E2E Admin Test Cobrand', 'E2E Manager Test Cobrand'];
	for (const name of testNames) {
		const found = await backendClient.collection('cobrands').getFullList({ filter: `name = "${name}"` });
		for (const cobrand of found) {
			const admins = await backendClient
				.collection('cobrandAdmins')
				.getFullList({ filter: `cobrand = "${c.id}"` });
			for (const admin of admins) await backendClient.collection('cobrandAdmins').delete(admin.id);
			const managers = await backendClient
				.collection('cobrandPropertyManagers')
				.getFullList({ filter: `cobrand = "${c.id}"` });
			for (const manager of managers) await backendClient.collection('cobrandPropertyManagers').delete(manager.id);
			await backendClient.collection('cobrands').delete(c.id);
		}
	}
}

test.beforeAll(async () => {
	const backendClient = await adminAuth();
	for (const [email, password] of [
		[TEST_EMAIL, TEST_PASSWORD],
		[TEST_EMAIL_2, TEST_PASSWORD_2]
	]) {
		try {
			await backendClient.collection('users').create({
				email,
				password,
				passwordConfirm: password
			});
		} catch {
			// user already exists
		}
	}
	await cleanupTestCobrands(backendClient);
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

test('create cobrand', async ({ page }) => {
	const backendClient = await adminAuth();

	await page.goto('/user/cobrands');
	await page.waitForLoadState('networkidle');
	await page.getByRole('button', { name: 'Add Cobrand' }).click();
	await page.waitForURL(/\/user\/cobrands\/add/);
	await page.getByLabel('Cobrand Name').fill('E2E Test Cobrand');
	await page.getByRole('button', { name: 'Save' }).click();
	await page.waitForURL(/\/user\/cobrands(\/)?$/);

	const cobrands = await backendClient
		.collection('cobrands')
		.getFullList({ filter: 'name = "E2E Test Cobrand"' });
	expect(cobrands.length).toBe(1);

	await cleanupTestCobrands(backendClient);
});

test('add admin to cobrand', async ({ page }) => {
	const backendClient = await adminAuth();

	const users1 = await backendClient.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
	const users2 = await backendClient.collection('users').getFullList({ filter: `email = "${TEST_EMAIL_2}"` });
	const testUserId = users1[0].id;
	const secondUserId = users2[0].id;

	const cobrand = await backendClient
		.collection('cobrands')
		.create<{ id: string }>({ name: 'E2E Admin Test Cobrand' });
	// approved: true — inviting a co-admin requires the inviter to be an approved
	// admin themselves (Phase 12's approval gate).
	await backendClient
		.collection('cobrandAdmins')
		.create({ user: testUserId, cobrand: cobrand.id, approved: true });

	await page.goto(`/user/cobrands/detail?id=${cobrand.id}`);
	await page.waitForLoadState('networkidle');

	await page.getByLabel('User ID').fill(secondUserId);
	await page.getByRole('button', { name: 'Add Admin' }).click();
	await page.waitForTimeout(1000);

	const admins = await backendClient
		.collection('cobrandAdmins')
		.getFullList({ filter: `cobrand = "${cobrand.id}"` });
	expect(admins.some((admin) => admin.user === secondUserId)).toBe(true);

	await cleanupTestCobrands(backendClient);
});

test('add property manager to cobrand', async ({ page }) => {
	const backendClient = await adminAuth();

	const users1 = await backendClient.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
	const testUserId = users1[0].id;
	const personId = await ensureTestPerson(backendClient, testUserId);

	const cobrand = await backendClient
		.collection('cobrands')
		.create<{ id: string }>({ name: 'E2E Manager Test Cobrand' });
	await backendClient.collection('cobrandAdmins').create({ user: testUserId, cobrand: cobrand.id });

	const property = await backendClient
		.collection('properties')
		.create<{ id: string }>({ address: '77 E2E Manager Ave' });
	const propertyOwner = await backendClient
		.collection('propertyOwners')
		.create<{ id: string }>({ property: property.id });
	await backendClient
		.collection('personPropertyOwners')
		.create({ propertyOwner: propertyOwner.id, person: personId });

	await page.goto(`/user/cobrands/detail?id=${cobrand.id}`);
	await page.waitForLoadState('networkidle');

	await page.getByLabel('Property ID').fill(property.id);
	await page.getByRole('button', { name: 'Add Property Manager' }).click();
	await page.waitForTimeout(1000);

	const managers = await backendClient
		.collection('cobrandPropertyManagers')
		.getFullList({ filter: `cobrand = "${cobrand.id}"` });
	expect(managers.some((manager) => manager.property === property.id)).toBe(true);

	for (const manager of managers) await backendClient.collection('cobrandPropertyManagers').delete(manager.id);
	await cleanupTestCobrands(backendClient);
	// Skip propertyOwner/property — the last-owner hook guard blocks deletion.
	const personPropertyOwners = await backendClient
		.collection('personPropertyOwners')
		.getFullList({ filter: `propertyOwner = "${propertyOwner.id}"` });
	for (const personPropertyOwner of personPropertyOwners) await backendClient.collection('personPropertyOwners').delete(personPropertyOwner.id);
});
