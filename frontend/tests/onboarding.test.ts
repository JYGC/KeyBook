import { test, expect, type Page } from '@playwright/test';
import PocketBase from 'pocketbase';

const PB_URL = 'http://192.168.8.144:8090';
const ADMIN_EMAIL = 'casperchen91@hotmail.com';
const ADMIN_PASSWORD = 'w3m#@tpth100';
const TEST_PASSWORD = 'Onboarding_test1';

const MONTH_NAMES = [
	'January',
	'February',
	'March',
	'April',
	'May',
	'June',
	'July',
	'August',
	'September',
	'October',
	'November',
	'December'
];

// The Date of Birth field is a read-only calendar picker (see full-collections
// PersonEditor DOB fix) — dates must be selected via the flatpickr calendar UI
// rather than typed, so tests navigate the calendar instead of using .fill().
async function pickDate(page: Page, labelText: string, isoDate: string): Promise<void> {
	const [year, month, day] = isoDate.split('-').map(Number);
	await page.getByLabel(labelText).click();
	const calendar = page.locator('.flatpickr-calendar.open');
	await calendar.waitFor({ state: 'visible' });

	const yearInput = calendar.locator('input.cur-year');
	await yearInput.fill(String(year));
	await yearInput.press('Enter');

	for (let monthStepAttempt = 0; monthStepAttempt < 24; monthStepAttempt++) {
		const currentMonthName = (await calendar.locator('.cur-month').textContent())?.trim();
		const currentYear = Number(await yearInput.inputValue());
		if (currentMonthName === MONTH_NAMES[month - 1] && currentYear === year) break;
		const currentIndex = MONTH_NAMES.indexOf(currentMonthName ?? '') + currentYear * 12;
		const targetIndex = month - 1 + year * 12;
		await calendar
			.locator(targetIndex < currentIndex ? '.flatpickr-prev-month' : '.flatpickr-next-month')
			.click();
	}

	await calendar
		.locator('.flatpickr-day:not(.prevMonthDay):not(.nextMonthDay)', {
			hasText: new RegExp(`^${day}$`)
		})
		.click();
}

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

async function ensureUser(
	backendClient: PocketBase,
	email: string,
	password: string
): Promise<string> {
	const existing = await backendClient
		.collection('users')
		.getFullList({ filter: `email = "${email}"` });
	if (existing.length > 0) return existing[0].id;
	const user = await backendClient
		.collection('users')
		.create<{ id: string }>({ email, password, passwordConfirm: password });
	return user.id;
}

async function cleanupUserByEmail(backendClient: PocketBase, email: string): Promise<void> {
	const users = await backendClient
		.collection('users')
		.getFullList({ filter: `email = "${email}"` });
	for (const user of users) await backendClient.collection('users').delete(user.id);
}

async function cleanupPersonsFor(backendClient: PocketBase, userId: string): Promise<void> {
	const persons = await backendClient
		.collection('persons')
		.getFullList({ filter: `user = "${userId}"` });
	for (const person of persons) await backendClient.collection('persons').delete(person.id);
}

async function cleanupCobrandsNamed(backendClient: PocketBase, name: string): Promise<void> {
	const found = await backendClient
		.collection('cobrands')
		.getFullList({ filter: `name = "${name}"` });
	for (const cobrand of found) {
		const admins = await backendClient
			.collection('cobrandAdmins')
			.getFullList({ filter: `cobrand = "${cobrand.id}"` });
		for (const admin of admins) await backendClient.collection('cobrandAdmins').delete(admin.id);
		await backendClient.collection('cobrands').delete(cobrand.id);
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

	const backendClient = await adminAuth();
	await cleanupUserByEmail(backendClient, email);
});

test('choosing person setup creates a linked person and lands on the property list', async ({
	page
}) => {
	const email = `e2e_onboarding_person_${Date.now()}@keybook.test`;
	const backendClient = await adminAuth();
	const userId = await ensureUser(backendClient, email, TEST_PASSWORD);

	await page.goto('/auth/login');
	await page.getByLabel('Email').fill(email);
	await page.getByLabel('Password').fill(TEST_PASSWORD);
	await page.getByText('Log in!').click();
	await page.waitForURL(/\/user\/setup/);

	await page.getByRole('button', { name: 'Set up my person profile' }).click();
	await page.waitForURL(/\/user\/persons\/setup/);
	await page.getByLabel('Name').fill('E2E Onboarding Person');
	await pickDate(page, 'Date of Birth', '1992-04-10');
	await page.getByRole('button', { name: 'Save' }).click();

	await page.waitForURL(/\/user\/properties\/list/);

	const persons = await backendClient
		.collection('persons')
		.getFullList({ filter: `user = "${userId}"` });
	expect(persons.length).toBe(1);
	expect(persons[0].name).toBe('E2E Onboarding Person');

	await cleanupPersonsFor(backendClient, userId);
	await cleanupUserByEmail(backendClient, email);
});

test('choosing cobrand setup creates the cobrand and shows the pending-approval notice', async ({
	page
}) => {
	const email = `e2e_onboarding_cobrand_${Date.now()}@keybook.test`;
	const cobrandName = `E2E Onboarding Cobrand ${Date.now()}`;
	const backendClient = await adminAuth();
	await ensureUser(backendClient, email, TEST_PASSWORD);

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

	const cobrands = await backendClient
		.collection('cobrands')
		.getFullList({ filter: `name = "${cobrandName}"` });
	expect(cobrands.length).toBe(1);
	const cobrandId = cobrands[0].id;

	const admins = await backendClient
		.collection('cobrandAdmins')
		.getFullList({ filter: `cobrand = "${cobrandId}"` });
	expect(admins.length).toBe(1);
	expect(admins[0].approved).toBe(false);

	await page.goto(`/user/cobrands/detail?id=${cobrandId}`);
	await page.waitForLoadState('networkidle');
	await expect(page.getByText('Pending approval')).toBeVisible();

	await cleanupCobrandsNamed(backendClient, cobrandName);
	await cleanupUserByEmail(backendClient, email);
});

test('a user with a linked person always lands on the property list from /user, /user/setup, or /user/persons/setup', async ({
	page
}) => {
	const email = `e2e_onboarding_haveperson_${Date.now()}@keybook.test`;
	const backendClient = await adminAuth();
	const userId = await ensureUser(backendClient, email, TEST_PASSWORD);
	await backendClient.collection('persons').create({
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

	await cleanupPersonsFor(backendClient, userId);
	await cleanupUserByEmail(backendClient, email);
});

test('a user with only a linked cobrandAdmins record lands on the cobrand list from /user or /user/setup', async ({
	page
}) => {
	const email = `e2e_onboarding_haveadmin_${Date.now()}@keybook.test`;
	const cobrandName = `E2E Onboarding Existing Cobrand ${Date.now()}`;
	const backendClient = await adminAuth();
	const userId = await ensureUser(backendClient, email, TEST_PASSWORD);
	const cobrand = await backendClient
		.collection('cobrands')
		.create<{ id: string }>({ name: cobrandName });
	await backendClient.collection('cobrandAdmins').create({ user: userId, cobrand: cobrand.id });

	await page.goto('/auth/login');
	await page.getByLabel('Email').fill(email);
	await page.getByLabel('Password').fill(TEST_PASSWORD);
	await page.getByText('Log in!').click();
	await page.waitForURL(/\/user\/cobrands/);

	await page.goto('/user/setup');
	await page.waitForURL(/\/user\/cobrands/);

	await cleanupCobrandsNamed(backendClient, cobrandName);
	await cleanupUserByEmail(backendClient, email);
});
