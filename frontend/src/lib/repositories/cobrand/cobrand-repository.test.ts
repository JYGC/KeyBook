import PocketBase from 'pocketbase';
import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { CobrandRepository } from './cobrand-repository';
import { CobrandAdminRepository } from './cobrand-admin-repository';
import { CobrandPropertyManagerRepository } from './cobrand-property-manager-repository';
import { CobrandPropertyOwnerRepository } from './cobrand-property-owner-repository';

const PB_URL = 'http://192.168.8.144:8090';
const backendClient = new PocketBase(PB_URL);

const INT_USER_EMAIL = 'cobrand_int_user@keybook.test';
const INT_USER_PASSWORD = 'CobrandInt_user1';

beforeAll(async () => {
	// PocketBase Go v0.22 uses /api/admins (not _superusers which is v0.23+).
	const adminAuthResponse = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ identity: 'casperchen91@hotmail.com', password: 'w3m#@tpth100' })
	});
	const adminAuthResult = await adminAuthResponse.json();
	backendClient.authStore.save(adminAuthResult.token, adminAuthResult.admin);

	const existing = await backendClient
		.collection('users')
		.getFullList({ filter: `email = "${INT_USER_EMAIL}"` });
	if (existing.length === 0) {
		await backendClient.collection('users').create({
			email: INT_USER_EMAIL,
			password: INT_USER_PASSWORD,
			passwordConfirm: INT_USER_PASSWORD
		});
	}
});

afterAll(async () => {
	const users = await backendClient.collection('users').getFullList({ filter: `email = "${INT_USER_EMAIL}"` });
	for (const user of users) await backendClient.collection('users').delete(user.id);
});

describe('CobrandRepository', () => {
	it('creates, reads, updates, and deletes a cobrand', async () => {
		const cobrandRepository = new CobrandRepository(backendClient);

		const created = await cobrandRepository.createCobrand('Test Firm');
		expect(created.name).toBe('Test Firm');
		expect(created.id).toBeTruthy();

		const fetched = await cobrandRepository.getCobrandById(created.id);
		expect(fetched.id).toBe(created.id);
		expect(fetched.name).toBe('Test Firm');

		await cobrandRepository.updateCobrand(created.id, 'Updated Firm');
		const updated = await cobrandRepository.getCobrandById(created.id);
		expect(updated.name).toBe('Updated Firm');

		const allCobrands = await cobrandRepository.getAllCobrands();
		expect(allCobrands.some((cobrand) => cobrand.id === created.id)).toBe(true);

		await cobrandRepository.deleteCobrand(created.id);
	});
});

describe('CobrandAdminRepository', () => {
	it('creates, queries, and deletes a cobrand admin', async () => {
		const cobrandRepository = new CobrandRepository(backendClient);
		const cobrandAdminRepository = new CobrandAdminRepository(backendClient);

		const cobrand = await cobrandRepository.createCobrand('Admin Test Firm');
		const users = await backendClient
			.collection('users')
			.getFullList({ filter: `email = "${INT_USER_EMAIL}"` });
		const userId = users[0].id;

		const created = await cobrandAdminRepository.createCobrandAdmin(userId, cobrand.id);
		expect(created.user).toBe(userId);
		expect(created.cobrand).toBe(cobrand.id);

		const byCobrand = await cobrandAdminRepository.getCobrandAdminsByCobrandId(cobrand.id);
		expect(byCobrand.some((cobrandAdmin) => cobrandAdmin.id === created.id)).toBe(true);

		const byUser = await cobrandAdminRepository.getCobrandAdminByUserId(userId);
		expect(byUser?.id).toBe(created.id);

		await cobrandAdminRepository.deleteCobrandAdmin(created.id);
		await cobrandRepository.deleteCobrand(cobrand.id);
	});

	it('returns null from getByUserId when the user has no admin record', async () => {
		const cobrandAdminRepository = new CobrandAdminRepository(backendClient);
		const result = await cobrandAdminRepository.getCobrandAdminByUserId('nonexistent-user-id');
		expect(result).toBeNull();
	});
});

describe('CobrandPropertyManagerRepository', () => {
	it('creates, queries, and deletes a cobrand property manager', async () => {
		const cobrandRepository = new CobrandRepository(backendClient);
		const cobrandPropertyManagerRepository = new CobrandPropertyManagerRepository(backendClient);

		const cobrand = await cobrandRepository.createCobrand('Manager Test Firm');
		const property = await backendClient.collection('properties').create<{ id: string }>({
			address: '99 Manager Test Ave'
		});

		const created = await cobrandPropertyManagerRepository.createCobrandPropertyManager(cobrand.id, property.id);
		expect(created.cobrand).toBe(cobrand.id);
		expect(created.property).toBe(property.id);

		const byCobrand = await cobrandPropertyManagerRepository.getCobrandPropertyManagersByCobrandId(cobrand.id);
		expect(byCobrand.some((cobrandPropertyManager) => cobrandPropertyManager.id === created.id)).toBe(true);

		const byProperty = await cobrandPropertyManagerRepository.getCobrandPropertyManagersByPropertyId(property.id);
		expect(byProperty.some((cobrandPropertyManager) => cobrandPropertyManager.id === created.id)).toBe(true);

		await cobrandPropertyManagerRepository.deleteCobrandPropertyManager(created.id);
		await backendClient.collection('properties').delete(property.id);
		await cobrandRepository.deleteCobrand(cobrand.id);
	});
});

describe('CobrandPropertyOwnerRepository', () => {
	it('creates, queries, and deletes a cobrand property owner', async () => {
		const cobrandRepository = new CobrandRepository(backendClient);
		const cobrandPropertyOwnerRepository = new CobrandPropertyOwnerRepository(backendClient);

		const cobrand = await cobrandRepository.createCobrand('Owner Test Firm');
		const property = await backendClient.collection('properties').create<{ id: string }>({
			address: '88 Owner Test Ave'
		});
		const propertyOwner = await backendClient.collection('propertyOwners').create<{ id: string }>({
			property: property.id
		});

		const created = await cobrandPropertyOwnerRepository.createCobrandPropertyOwner(cobrand.id, propertyOwner.id);
		expect(created.cobrand).toBe(cobrand.id);
		expect(created.propertyOwner).toBe(propertyOwner.id);

		const byCobrand = await cobrandPropertyOwnerRepository.getCobrandPropertyOwnersByCobrandId(cobrand.id);
		expect(byCobrand.some((cobrandPropertyOwner) => cobrandPropertyOwner.id === created.id)).toBe(true);

		const byPropertyOwner = await cobrandPropertyOwnerRepository.getCobrandPropertyOwnersByPropertyOwnerId(propertyOwner.id);
		expect(byPropertyOwner.some((cobrandPropertyOwner) => cobrandPropertyOwner.id === created.id)).toBe(true);

		await cobrandPropertyOwnerRepository.deleteCobrandPropertyOwner(created.id);
		// propertyOwner cannot be deleted when it is the last owner (backend hook guard),
		// and property cannot be deleted while propertyOwner still references it.
		// Leave these small orphan records; they do not affect other tests.
		await cobrandRepository.deleteCobrand(cobrand.id);
	});
});
