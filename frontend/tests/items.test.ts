import { test, expect } from '@playwright/test';
import PocketBase from 'pocketbase';

const PB_URL = 'http://192.168.8.144:8090';
const ADMIN_EMAIL = 'casperchen91@hotmail.com';
const ADMIN_PASSWORD = 'w3m#@tpth100';
const TEST_EMAIL = 'e2e_items@keybook.test';
const TEST_PASSWORD = 'E2Eitems_test1';

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

async function ensureTestPerson(backendClient: PocketBase): Promise<string> {
	const users = await backendClient
		.collection('users')
		.getFullList({ filter: `email = "${TEST_EMAIL}"` });
	const testUser = users[0];
	const persons = await backendClient
		.collection('persons')
		.getFullList({ filter: `user = "${testUser.id}"` });
	if (persons.length > 0) return persons[0].id;
	const person = await backendClient.collection('persons').create({
		name: 'E2E Test Person',
		DOB: '1990-01-01',
		user: testUser.id
	});
	return person.id;
}

async function cleanupTestUserItems(backendClient: PocketBase): Promise<void> {
	const users = await backendClient
		.collection('users')
		.getFullList({ filter: `email = "${TEST_EMAIL}"` });
	if (users.length === 0) return;
	const persons = await backendClient
		.collection('persons')
		.getFullList({ filter: `user = "${users[0].id}"` });
	for (const person of persons) {
		const personItems = await backendClient
			.collection('personItems')
			.getFullList({ filter: `person = "${person.id}"` });
		for (const personItem of personItems) {
			const itemId = personItem.item as string;
			try {
				const entryDevices = await backendClient
					.collection('entryDevices')
					.getFullList({ filter: `item = "${itemId}"` });
				for (const entryDevice of entryDevices)
					await backendClient.collection('entryDevices').delete(entryDevice.id);
			} catch {
				/* ignore */
			}
			try {
				await backendClient.collection('personItems').delete(personItem.id);
			} catch {
				/* ignore */
			}
			try {
				await backendClient.collection('items').delete(itemId);
			} catch {
				/* ignore */
			}
		}
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
	await cleanupTestUserItems(backendClient);
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

test('create item', async ({ page }) => {
	const backendClient = await adminAuth();

	await page.goto('/user/items');
	await page.waitForLoadState('networkidle');
	await page.getByRole('button', { name: 'Add Item' }).click();
	await page.waitForURL(/\/user\/items\/add/);
	await page.getByLabel('Item Name').fill('Front Door Key');
	await page.getByLabel('Description').fill('Key to the front door');
	await page.getByRole('button', { name: 'Save' }).click();
	await page.waitForURL(/\/user\/items(\/)?$/);

	const items = await backendClient
		.collection('items')
		.getFullList({ filter: 'name = "Front Door Key"' });
	expect(items.length).toBe(1);
	expect(items[0].description).toBe('Key to the front door');
	await backendClient.collection('items').delete(items[0].id);
});

test('edit item name', async ({ page }) => {
	const backendClient = await adminAuth();
	const personId = await ensureTestPerson(backendClient);
	const item = await backendClient
		.collection('items')
		.create({ name: 'E2E Edit Item', description: 'Original' });
	await backendClient.collection('personItems').create({ person: personId, item: item.id });

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

	const items = await backendClient
		.collection('items')
		.getFullList({ filter: 'name = "E2E Edited Key"' });
	await backendClient
		.collection('personItems')
		.delete(
			(
				await backendClient
					.collection('personItems')
					.getFullList({ filter: `item = "${items[0].id}"` })
			)[0].id
		);
	await backendClient.collection('items').delete(items[0].id);
});

test('delete item', async ({ page }) => {
	const backendClient = await adminAuth();
	const personId = await ensureTestPerson(backendClient);
	const item = await backendClient
		.collection('items')
		.create({ name: 'E2E Delete Item', description: 'Gone soon' });
	await backendClient.collection('personItems').create({ person: personId, item: item.id });

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
	const backendClient = await adminAuth();
	const personId = await ensureTestPerson(backendClient);
	const item = await backendClient
		.collection('items')
		.create({ name: 'E2E Key With Device', description: 'A test key' });
	await backendClient.collection('personItems').create({ person: personId, item: item.id });

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

	const entryDevices = await backendClient
		.collection('entryDevices')
		.getFullList({ filter: `item = "${item.id}"` });
	expect(entryDevices.length).toBe(1);
	expect(entryDevices[0].identifier).toBe('K-001');

	await backendClient.collection('entryDevices').delete(entryDevices[0].id);
	await backendClient
		.collection('personItems')
		.delete(
			(
				await backendClient.collection('personItems').getFullList({ filter: `item = "${item.id}"` })
			)[0].id
		);
	await backendClient.collection('items').delete(item.id);
});

test('mark entry device defunct', async ({ page }) => {
	const backendClient = await adminAuth();
	const personId = await ensureTestPerson(backendClient);
	const item = await backendClient
		.collection('items')
		.create({ name: 'E2E Defunct Key', description: 'A key going defunct' });
	await backendClient.collection('personItems').create({ person: personId, item: item.id });
	const entryDevice = await backendClient.collection('entryDevices').create({
		item: item.id,
		deviceType: 'Key',
		identifier: 'K-002',
		defunctReason: 'None'
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

	const updated = await backendClient.collection('entryDevices').getOne(entryDevice.id);
	expect(updated.defunctReason).toBe('Lost');

	await backendClient.collection('entryDevices').delete(entryDevice.id);
	await backendClient
		.collection('personItems')
		.delete(
			(
				await backendClient.collection('personItems').getFullList({ filter: `item = "${item.id}"` })
			)[0].id
		);
	await backendClient.collection('items').delete(item.id);
});
