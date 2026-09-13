import PocketBase from 'pocketbase';
import { describe, it, expect, beforeAll } from 'vitest';
import { PersonRepository } from './person-repository';
import { PersonPropertyOwnerRepository } from './person-property-owner-repository';

const PB_URL = 'http://192.168.8.144:8090';
const backendClient = new PocketBase(PB_URL);

beforeAll(async () => {
	const adminAuthResponse = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ identity: 'casperchen91@hotmail.com', password: 'w3m#@tpth100' })
	});
	const adminAuthResult = await adminAuthResponse.json();
	backendClient.authStore.save(adminAuthResult.token, adminAuthResult.admin);
});

describe('PersonRepository', () => {
	it('creates, reads, updates, and deletes a person', async () => {
		const personRepository = new PersonRepository(backendClient);

		const created = await personRepository.createPerson('Person Repo Test', '1990-06-15');
		expect(created.name).toBe('Person Repo Test');
		expect(created.DOB).toMatch(/1990-06-15/);
		expect(created.id).toBeTruthy();

		const fetched = await personRepository.getPersonById(created.id);
		expect(fetched.id).toBe(created.id);
		expect(fetched.name).toBe('Person Repo Test');

		const allPersons = await personRepository.getAllPersons();
		expect(allPersons.some((person) => person.id === created.id)).toBe(true);

		const updated = await personRepository.updatePerson(created.id, 'Person Repo Test Updated', '1990-06-15');
		expect(updated.name).toBe('Person Repo Test Updated');

		await personRepository.deletePerson(created.id);
		const allPersonsAfterDelete = await personRepository.getAllPersons();
		expect(allPersonsAfterDelete.some((personsAfterDelete) => personsAfterDelete.id === created.id)).toBe(false);
	});

	it('links the person to a user when a userId is provided', async () => {
		const personRepository = new PersonRepository(backendClient);
		const user = await backendClient.collection('users').create<{ id: string }>({
			email: `person-repo-link-${Date.now()}@keybook.test`,
			password: 'PersonRepoLink1!',
			passwordConfirm: 'PersonRepoLink1!'
		});

		const created = await personRepository.createPerson('Person Repo Link Test', '1990-06-15', user.id);
		expect(created.user).toBe(user.id);

		const fetched = await personRepository.getPersonByUserId(user.id);
		expect(fetched?.id).toBe(created.id);

		await personRepository.deletePerson(created.id);
		await backendClient.collection('users').delete(user.id);
	});
});

describe('PersonPropertyOwnerRepository', () => {
	it('creates, queries, and deletes a person-property-owner link', async () => {
		const personRepository = new PersonRepository(backendClient);
		const personPropertyOwnerRepository = new PersonPropertyOwnerRepository(backendClient);

		const person = await personRepository.createPerson('PPO Repo Test Person', '1985-03-20');
		const property = await backendClient
			.collection('properties')
			.create<{ id: string }>({ address: '10 PPO Repo Test Ave' });
		const propertyOwner = await backendClient
			.collection('propertyOwners')
			.create<{ id: string }>({ property: property.id });

		const created = await personPropertyOwnerRepository.createPersonPropertyOwner(person.id, propertyOwner.id);
		expect(created.person).toBe(person.id);
		expect(created.propertyOwner).toBe(propertyOwner.id);
		expect(created.id).toBeTruthy();

		const byPerson = await personPropertyOwnerRepository.getPersonPropertyOwnersByPersonId(person.id);
		expect(byPerson.some((personPropertyOwner) => personPropertyOwner.id === created.id)).toBe(true);

		const byPropertyOwner = await personPropertyOwnerRepository.getPersonPropertyOwnersByPropertyOwnerId(propertyOwner.id);
		expect(byPropertyOwner.some((personPropertyOwner) => personPropertyOwner.id === created.id)).toBe(true);

		await personPropertyOwnerRepository.deletePersonPropertyOwner(created.id);
		await personRepository.deletePerson(person.id);
		// Skip propertyOwner and property deletion — last-owner hook guard prevents it
	});
});
