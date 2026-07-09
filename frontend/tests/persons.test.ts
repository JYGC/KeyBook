import { test, expect, type Page } from '@playwright/test';
import PocketBase from 'pocketbase';

const PB_URL = 'http://192.168.8.144:8090';
const ADMIN_EMAIL = 'casperchen91@hotmail.com';
const ADMIN_PASSWORD = 'w3m#@tpth100';
const TEST_EMAIL = 'e2e_persons@keybook.test';
const TEST_PASSWORD = 'E2Epersons_test1';

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

	for (let i = 0; i < 24; i++) {
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

async function cleanupTestPersons(pb: PocketBase): Promise<void> {
	const persons = await pb.collection('persons').getFullList({
		filter: 'name ~ "E2E Persons Test"'
	});
	for (const p of persons) {
		const ppos = await pb.collection('personPropertyOwners').getFullList({
			filter: `person = "${p.id}"`
		});
		for (const ppo of ppos) await pb.collection('personPropertyOwners').delete(ppo.id);
		await pb.collection('persons').delete(p.id);
	}
}

test.beforeAll(async () => {
	const pb = await adminAuth();
	try {
		await pb.collection('users').create({
			email: TEST_EMAIL,
			password: TEST_PASSWORD,
			passwordConfirm: TEST_PASSWORD
		});
	} catch {
		// user already exists
	}
	await cleanupTestPersons(pb);
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

test('create person', async ({ page }) => {
	const pb = await adminAuth();

	await page.goto('/user/persons');
	await page.waitForLoadState('networkidle');
	await page.getByRole('button', { name: 'Add Person' }).click();
	await page.waitForURL(/\/user\/persons\/add/);

	await page.getByLabel('Name').fill('E2E Persons Test Person');
	await pickDate(page, 'Date of Birth', '1990-06-15');
	await page.getByRole('button', { name: 'Save' }).click();
	await page.waitForURL(/\/user\/persons(\/)?$/);

	const persons = await pb.collection('persons').getFullList({
		filter: 'name = "E2E Persons Test Person"'
	});
	expect(persons.length).toBe(1);
	expect(persons[0].DOB).toBe('1990-06-15 00:00:00.000Z');

	await cleanupTestPersons(pb);
});

test('edit person and view roles', async ({ page }) => {
	const pb = await adminAuth();

	const users = await pb.collection('users').getFullList({ filter: `email = "${TEST_EMAIL}"` });
	const testUserId = users[0].id;
	const person = await pb.collection('persons').create<{ id: string }>({
		name: 'E2E Persons Test Edit',
		DOB: '1985-03-20',
		user: testUserId
	});

	await page.goto(`/user/persons/detail?id=${person.id}`);
	await page.waitForLoadState('networkidle');

	await page.getByLabel('Name').fill('E2E Persons Test Edited');
	await page.getByRole('button', { name: 'Save' }).click();
	await page.waitForURL(/\/user\/persons(\/)?$/);

	const updated = await pb.collection('persons').getOne<{ name: string }>(person.id);
	expect(updated.name).toBe('E2E Persons Test Edited');

	await cleanupTestPersons(pb);
});
