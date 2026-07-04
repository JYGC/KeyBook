import PocketBase from 'pocketbase';
import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { CobrandRepository } from './cobrand-repository';
import { CobrandAdminRepository } from './cobrand-admin-repository';
import { CobrandPropertyManagerRepository } from './cobrand-property-manager-repository';
import { CobrandPropertyOwnerRepository } from './cobrand-property-owner-repository';

const PB_URL = 'http://192.168.8.144:8090';
const pb = new PocketBase(PB_URL);

const INT_USER_EMAIL = 'cobrand_int_user@keybook.test';
const INT_USER_PASSWORD = 'CobrandInt_user1';

beforeAll(async () => {
	// PocketBase Go v0.22 uses /api/admins (not _superusers which is v0.23+).
	const res = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ identity: 'casperchen91@hotmail.com', password: 'w3m#@tpth100' })
	});
	const data = await res.json();
	pb.authStore.save(data.token, data.admin);

	// Create integration test user if not exists
	const existing = await pb
		.collection('users')
		.getFullList({ filter: `email = "${INT_USER_EMAIL}"` });
	if (existing.length === 0) {
		await pb.collection('users').create({
			email: INT_USER_EMAIL,
			password: INT_USER_PASSWORD,
			passwordConfirm: INT_USER_PASSWORD
		});
	}
});

afterAll(async () => {
	// Clean up integration test user
	const users = await pb.collection('users').getFullList({ filter: `email = "${INT_USER_EMAIL}"` });
	for (const u of users) await pb.collection('users').delete(u.id);
});

describe('CobrandRepository', () => {
	it('creates, reads, updates, and deletes a cobrand', async () => {
		const repo = new CobrandRepository(pb);

		const created = await repo.create('Test Firm');
		expect(created.name).toBe('Test Firm');
		expect(created.id).toBeTruthy();

		const fetched = await repo.getById(created.id);
		expect(fetched.id).toBe(created.id);
		expect(fetched.name).toBe('Test Firm');

		await repo.update(created.id, 'Updated Firm');
		const updated = await repo.getById(created.id);
		expect(updated.name).toBe('Updated Firm');

		const all = await repo.getAll();
		expect(all.some((c) => c.id === created.id)).toBe(true);

		await repo.delete(created.id);
	});
});

describe('CobrandAdminRepository', () => {
	it('creates, queries, and deletes a cobrand admin', async () => {
		const cobrandRepo = new CobrandRepository(pb);
		const adminRepo = new CobrandAdminRepository(pb);

		const cobrand = await cobrandRepo.create('Admin Test Firm');
		const users = await pb
			.collection('users')
			.getFullList({ filter: `email = "${INT_USER_EMAIL}"` });
		const userId = users[0].id;

		const created = await adminRepo.create(userId, cobrand.id);
		expect(created.user).toBe(userId);
		expect(created.cobrand).toBe(cobrand.id);

		const byCobrand = await adminRepo.getByCobrandId(cobrand.id);
		expect(byCobrand.some((a) => a.id === created.id)).toBe(true);

		const byUser = await adminRepo.getByUserId(userId);
		expect(byUser?.id).toBe(created.id);

		await adminRepo.delete(created.id);
		await cobrandRepo.delete(cobrand.id);
	});

	it('returns null from getByUserId when the user has no admin record', async () => {
		const adminRepo = new CobrandAdminRepository(pb);
		const result = await adminRepo.getByUserId('nonexistent-user-id');
		expect(result).toBeNull();
	});
});

describe('CobrandPropertyManagerRepository', () => {
	it('creates, queries, and deletes a cobrand property manager', async () => {
		const cobrandRepo = new CobrandRepository(pb);
		const managerRepo = new CobrandPropertyManagerRepository(pb);

		const cobrand = await cobrandRepo.create('Manager Test Firm');
		const property = await pb.collection('properties').create<{ id: string }>({
			address: '99 Manager Test Ave'
		});

		const created = await managerRepo.create(cobrand.id, property.id);
		expect(created.cobrand).toBe(cobrand.id);
		expect(created.property).toBe(property.id);

		const byCobrand = await managerRepo.getByCobrandId(cobrand.id);
		expect(byCobrand.some((m) => m.id === created.id)).toBe(true);

		const byProperty = await managerRepo.getByPropertyId(property.id);
		expect(byProperty.some((m) => m.id === created.id)).toBe(true);

		await managerRepo.delete(created.id);
		await pb.collection('properties').delete(property.id);
		await cobrandRepo.delete(cobrand.id);
	});
});

describe('CobrandPropertyOwnerRepository', () => {
	it('creates, queries, and deletes a cobrand property owner', async () => {
		const cobrandRepo = new CobrandRepository(pb);
		const ownerRepo = new CobrandPropertyOwnerRepository(pb);

		const cobrand = await cobrandRepo.create('Owner Test Firm');
		const property = await pb.collection('properties').create<{ id: string }>({
			address: '88 Owner Test Ave'
		});
		const propertyOwner = await pb.collection('propertyOwners').create<{ id: string }>({
			property: property.id
		});

		const created = await ownerRepo.create(cobrand.id, propertyOwner.id);
		expect(created.cobrand).toBe(cobrand.id);
		expect(created.propertyOwner).toBe(propertyOwner.id);

		const byCobrand = await ownerRepo.getByCobrandId(cobrand.id);
		expect(byCobrand.some((o) => o.id === created.id)).toBe(true);

		const byPropertyOwner = await ownerRepo.getByPropertyOwnerId(propertyOwner.id);
		expect(byPropertyOwner.some((o) => o.id === created.id)).toBe(true);

		await ownerRepo.delete(created.id);
		// propertyOwner cannot be deleted when it is the last owner (backend hook guard),
		// and property cannot be deleted while propertyOwner still references it.
		// Leave these small orphan records; they do not affect other tests.
		await cobrandRepo.delete(cobrand.id);
	});
});
