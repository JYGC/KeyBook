import PocketBase from 'pocketbase';
import type { IPersonModel } from '$lib/models/person-models';

export interface IPersonRepository {
	getAllPersons(): Promise<IPersonModel[]>;
	getPersonById(id: string): Promise<IPersonModel>;
	getPersonByUserId(userId: string): Promise<IPersonModel | null>;
	createPerson(name: string, dob: string, userId?: string): Promise<IPersonModel>;
	updatePerson(id: string, name: string, dob: string): Promise<IPersonModel>;
	deletePerson(id: string): Promise<void>;
}

export class PersonRepository implements IPersonRepository {
	constructor(private readonly backendClient: PocketBase) {}

	async getAllPersons(): Promise<IPersonModel[]> {
		return await this.backendClient.collection('persons').getFullList<IPersonModel>();
	}

	async getPersonById(id: string): Promise<IPersonModel> {
		const allPersons = await this.backendClient.collection('persons').getFullList<IPersonModel>();
		const person = allPersons.find((candidatePerson) => candidatePerson.id === id);
		if (!person) throw new Error(`Person not found: ${id}`);
		return person;
	}

	async getPersonByUserId(userId: string): Promise<IPersonModel | null> {
		const allPersons = await this.backendClient.collection('persons').getFullList<IPersonModel>();
		return allPersons.find((person) => person.user === userId) ?? null;
	}

	async createPerson(name: string, dob: string, userId?: string): Promise<IPersonModel> {
		const payload: Record<string, string> = { name, DOB: dob };
		if (userId) payload.user = userId;
		return await this.backendClient.collection('persons').create<IPersonModel>(payload);
	}

	async updatePerson(id: string, name: string, dob: string): Promise<IPersonModel> {
		return await this.backendClient
			.collection('persons')
			.update<IPersonModel>(id, { name, DOB: dob });
	}

	async deletePerson(id: string): Promise<void> {
		await this.backendClient.collection('persons').delete(id);
	}
}
