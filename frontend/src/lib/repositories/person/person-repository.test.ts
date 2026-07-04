import PocketBase from 'pocketbase';
import { describe, it, expect, beforeAll } from 'vitest';
import { PersonRepository } from './person-repository';
import { PersonPropertyOwnerRepository } from './person-property-owner-repository';

const PB_URL = 'http://192.168.8.144:8090';
const pb = new PocketBase(PB_URL);

beforeAll(async () => {
	const res = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ identity: 'casperchen91@hotmail.com', password: 'w3m#@tpth100' })
	});
	const data = await res.json();
	pb.authStore.save(data.token, data.admin);
});

describe('PersonRepository', () => {
	it('creates, reads, updates, and deletes a person', async () => {
		const repo = new PersonRepository(pb);

		const created = await repo.create('Person Repo Test', '1990-06-15');
		expect(created.name).toBe('Person Repo Test');
		expect(created.DOB).toMatch(/1990-06-15/);
		expect(created.id).toBeTruthy();

		const fetched = await repo.getById(created.id);
		expect(fetched.id).toBe(created.id);
		expect(fetched.name).toBe('Person Repo Test');

		const all = await repo.getAll();
		expect(all.some((p) => p.id === created.id)).toBe(true);

		const updated = await repo.update(created.id, 'Person Repo Test Updated', '1990-06-15');
		expect(updated.name).toBe('Person Repo Test Updated');

		await repo.delete(created.id);
		const allAfter = await repo.getAll();
		expect(allAfter.some((p) => p.id === created.id)).toBe(false);
	});

	it('links the person to a user when a userId is provided', async () => {
		const repo = new PersonRepository(pb);
		const user = await pb.collection('users').create<{ id: string }>({
			email: `person-repo-link-${Date.now()}@keybook.test`,
			password: 'PersonRepoLink1!',
			passwordConfirm: 'PersonRepoLink1!'
		});

		const created = await repo.create('Person Repo Link Test', '1990-06-15', user.id);
		expect(created.user).toBe(user.id);

		const fetched = await repo.getByUserId(user.id);
		expect(fetched?.id).toBe(created.id);

		await repo.delete(created.id);
		await pb.collection('users').delete(user.id);
	});
});

describe('PersonPropertyOwnerRepository', () => {
	it('creates, queries, and deletes a person-property-owner link', async () => {
		const personRepo = new PersonRepository(pb);
		const ppoRepo = new PersonPropertyOwnerRepository(pb);

		const person = await personRepo.create('PPO Repo Test Person', '1985-03-20');
		const property = await pb
			.collection('properties')
			.create<{ id: string }>({ address: '10 PPO Repo Test Ave' });
		const propertyOwner = await pb
			.collection('propertyOwners')
			.create<{ id: string }>({ property: property.id });

		const created = await ppoRepo.create(person.id, propertyOwner.id);
		expect(created.person).toBe(person.id);
		expect(created.propertyOwner).toBe(propertyOwner.id);
		expect(created.id).toBeTruthy();

		const byPerson = await ppoRepo.getByPersonId(person.id);
		expect(byPerson.some((r) => r.id === created.id)).toBe(true);

		const byPropertyOwner = await ppoRepo.getByPropertyOwnerId(propertyOwner.id);
		expect(byPropertyOwner.some((r) => r.id === created.id)).toBe(true);

		await ppoRepo.delete(created.id);
		await personRepo.delete(person.id);
		// Skip propertyOwner and property deletion — last-owner hook guard prevents it
	});
});
