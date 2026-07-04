import { test, expect } from '@playwright/test';
import PocketBase from 'pocketbase';

const PB_URL = 'http://192.168.8.144:8090';
const ADMIN_EMAIL = 'casperchen91@hotmail.com';
const ADMIN_PASSWORD = 'w3m#@tpth100';
const TEST_PASSWORD = 'Onboarding_test1';

async function adminAuth(): Promise<PocketBase> {
	const pb = new PocketBase(PB_URL);
	const res = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ identity: ADMIN_EMAIL, password: ADMIN_PASSWORD })
	});
	const data = await res.json();
	pb.authStore.save(data.token, data.admin);
	return pb;
}

async function ensureUser(pb: PocketBase, email: string, password: string): Promise<string> {
	const existing = await pb.collection('users').getFullList({ filter: `email = "${email}"` });
	if (existing.length > 0) return existing[0].id;
	const user = await pb
		.collection('users')
		.create<{ id: string }>({ email, password, passwordConfirm: password });
	return user.id;
}

async function cleanupUserByEmail(pb: PocketBase, email: string): Promise<void> {
	const users = await pb.collection('users').getFullList({ filter: `email = "${email}"` });
	for (const u of users) await pb.collection('users').delete(u.id);
}

async function cleanupPersonsFor(pb: PocketBase, userId: string): Promise<void> {
	const persons = await pb.collection('persons').getFullList({ filter: `user = "${userId}"` });
	for (const p of persons) await pb.collection('persons').delete(p.id);
}

async function cleanupCobrandsNamed(pb: PocketBase, name: string): Promise<void> {
	const found = await pb.collection('cobrands').getFullList({ filter: `name = "${name}"` });
	for (const c of found) {
		const admins = await pb
			.collection('cobrandAdmins')
			.getFullList({ filter: `cobrand = "${c.id}"` });
		for (const a of admins) await pb.collection('cobrandAdmins').delete(a.id);
		await pb.collection('cobrands').delete(c.id);
	}
}

test('register lands on the account setup choice page', async ({ page }) => {
	const email = `e2e_onboarding_register_${Date.now()}@keybook.test`;

	await page.goto('/auth/register');
	await page.getByLabel('Email').fill(email);
	await page.getByLabel('Password', { exact: true }).fill(TEST_PASSWORD);
	await page.getByLabel('Confirm password', { exact: true }).fill(TEST_PASSWORD);
	await page.getByRole('button', { name: 'Create account!' }).click();

	await page.waitForURL(/\/user\/setup/);
	await expect(page.getByRole('button', { name: 'Set up my person profile' })).toBeVisible();
	await expect(page.getByRole('button', { name: 'Set up my company' })).toBeVisible();

	const pb = await adminAuth();
	await cleanupUserByEmail(pb, email);
});

test('choosing person setup creates a linked person and lands on the property list', async ({
	page
}) => {
	const email = `e2e_onboarding_person_${Date.now()}@keybook.test`;
	const pb = await adminAuth();
	const userId = await ensureUser(pb, email, TEST_PASSWORD);

	await page.goto('/auth/login');
	await page.getByLabel('Email').fill(email);
	await page.getByLabel('Password').fill(TEST_PASSWORD);
	await page.getByText('Log in!').click();
	await page.waitForURL(/\/user\/setup/);

	await page.getByRole('button', { name: 'Set up my person profile' }).click();
	await page.waitForURL(/\/user\/persons\/setup/);
	await page.getByLabel('Name').fill('E2E Onboarding Person');
	await page.getByLabel('Date of Birth').fill('1992-04-10');
	await page.getByRole('button', { name: 'Save' }).click();

	await page.waitForURL(/\/user\/properties\/list/);

	const persons = await pb.collection('persons').getFullList({ filter: `user = "${userId}"` });
	expect(persons.length).toBe(1);
	expect(persons[0].name).toBe('E2E Onboarding Person');

	await cleanupPersonsFor(pb, userId);
	await cleanupUserByEmail(pb, email);
});

test('choosing cobrand setup creates the cobrand and shows the pending-approval notice', async ({
	page
}) => {
	const email = `e2e_onboarding_cobrand_${Date.now()}@keybook.test`;
	const cobrandName = `E2E Onboarding Cobrand ${Date.now()}`;
	const pb = await adminAuth();
	await ensureUser(pb, email, TEST_PASSWORD);

	await page.goto('/auth/login');
	await page.getByLabel('Email').fill(email);
	await page.getByLabel('Password').fill(TEST_PASSWORD);
	await page.getByText('Log in!').click();
	await page.waitForURL(/\/user\/setup/);

	await page.getByRole('button', { name: 'Set up my company' }).click();
	await page.waitForURL(/\/user\/cobrands\/add/);
	await page.getByLabel('Cobrand Name').fill(cobrandName);
	await page.getByRole('button', { name: 'Save' }).click();

	// The existing CobrandAddModule navigates back via history.back(), which returns to
	// /user/setup; AccountSetupModule's own redirect-away guard then bounces it forward
	// to /user/cobrands/ once it sees the new admin record.
	await page.waitForURL(/\/user\/cobrands(\/)?$/, { timeout: 10000 });

	const cobrands = await pb
		.collection('cobrands')
		.getFullList({ filter: `name = "${cobrandName}"` });
	expect(cobrands.length).toBe(1);
	const cobrandId = cobrands[0].id;

	const admins = await pb
		.collection('cobrandAdmins')
		.getFullList({ filter: `cobrand = "${cobrandId}"` });
	expect(admins.length).toBe(1);
	expect(admins[0].approved).toBe(false);

	await page.goto(`/user/cobrands/detail?id=${cobrandId}`);
	await page.waitForLoadState('networkidle');
	await expect(page.getByText('Pending approval')).toBeVisible();

	await cleanupCobrandsNamed(pb, cobrandName);
	await cleanupUserByEmail(pb, email);
});

test('a user with a linked person always lands on the property list from /user, /user/setup, or /user/persons/setup', async ({
	page
}) => {
	const email = `e2e_onboarding_haveperson_${Date.now()}@keybook.test`;
	const pb = await adminAuth();
	const userId = await ensureUser(pb, email, TEST_PASSWORD);
	await pb.collection('persons').create({
		name: 'E2E Onboarding Existing Person',
		DOB: '1990-01-01',
		user: userId
	});

	await page.goto('/auth/login');
	await page.getByLabel('Email').fill(email);
	await page.getByLabel('Password').fill(TEST_PASSWORD);
	await page.getByText('Log in!').click();
	await page.waitForURL(/\/user\/properties\/list/);

	await page.goto('/user/setup');
	await page.waitForURL(/\/user\/properties\/list/);

	await page.goto('/user/persons/setup');
	await page.waitForURL(/\/user\/properties\/list/);

	await cleanupPersonsFor(pb, userId);
	await cleanupUserByEmail(pb, email);
});

test('a user with only a linked cobrandAdmins record lands on the cobrand list from /user or /user/setup', async ({
	page
}) => {
	const email = `e2e_onboarding_haveadmin_${Date.now()}@keybook.test`;
	const cobrandName = `E2E Onboarding Existing Cobrand ${Date.now()}`;
	const pb = await adminAuth();
	const userId = await ensureUser(pb, email, TEST_PASSWORD);
	const cobrand = await pb.collection('cobrands').create<{ id: string }>({ name: cobrandName });
	await pb.collection('cobrandAdmins').create({ user: userId, cobrand: cobrand.id });

	await page.goto('/auth/login');
	await page.getByLabel('Email').fill(email);
	await page.getByLabel('Password').fill(TEST_PASSWORD);
	await page.getByText('Log in!').click();
	await page.waitForURL(/\/user\/cobrands/);

	await page.goto('/user/setup');
	await page.waitForURL(/\/user\/cobrands/);

	await cleanupCobrandsNamed(pb, cobrandName);
	await cleanupUserByEmail(pb, email);
});
