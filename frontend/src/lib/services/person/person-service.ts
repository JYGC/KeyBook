import type { IPersonRepository } from '$lib/repositories/person/person-repository';
import type { IPersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
import type { ITenantRepository } from '$lib/repositories/property/tenant-repository';
import type { IHouseholdRepository } from '$lib/repositories/property/household-repository';
import type { IAgentRepository } from '$lib/repositories/agent/agent-repository';
import type { IPersonModel } from '$lib/models/person-models';

export interface IPersonService {
	validatePersonName(name: string): void;
	validatePersonDOB(dob: string): void;
	getAllPersons(): Promise<IPersonModel[]>;
	getPersonById(id: string): Promise<IPersonModel>;
	getPersonByUserId(userId: string): Promise<IPersonModel | null>;
	createPerson(name: string, dob: string, userId?: string): Promise<IPersonModel>;
	updatePerson(id: string, name: string, dob: string): Promise<IPersonModel>;
	deletePerson(id: string): Promise<void>;
	getRolesForPerson(personId: string): Promise<string[]>;
}

export class PersonService implements IPersonService {
	constructor(
		private readonly personRepository: IPersonRepository,
		private readonly personPropertyOwnerRepository: IPersonPropertyOwnerRepository,
		private readonly tenantRepository: ITenantRepository,
		private readonly householdRepository: IHouseholdRepository,
		private readonly agentRepository: IAgentRepository
	) {}

	validatePersonName(name: string): void {
		if (!name || name.trim() === '') throw new Error('person name is required');
	}

	validatePersonDOB(dob: string): void {
		if (!dob || dob.trim() === '') throw new Error('date of birth is required');
	}

	async getAllPersons(): Promise<IPersonModel[]> {
		return await this.personRepository.getAllPersons();
	}

	async getPersonById(id: string): Promise<IPersonModel> {
		return await this.personRepository.getPersonById(id);
	}

	async getPersonByUserId(userId: string): Promise<IPersonModel | null> {
		return await this.personRepository.getPersonByUserId(userId);
	}

	async createPerson(name: string, dob: string, userId?: string): Promise<IPersonModel> {
		this.validatePersonName(name);
		this.validatePersonDOB(dob);
		return await this.personRepository.createPerson(name, dob, userId);
	}

	async updatePerson(id: string, name: string, dob: string): Promise<IPersonModel> {
		this.validatePersonName(name);
		this.validatePersonDOB(dob);
		return await this.personRepository.updatePerson(id, name, dob);
	}

	async deletePerson(id: string): Promise<void> {
		await this.personRepository.deletePerson(id);
	}

	async getRolesForPerson(personId: string): Promise<string[]> {
		const roles: string[] = [];
		const [personPropertyOwners, tenants, households, agents] = await Promise.all([
			this.personPropertyOwnerRepository.getPersonPropertyOwnersByPersonId(personId),
			this.tenantRepository.getTenantsByPersonId(personId),
			this.householdRepository.getHouseholdsByPersonId(personId),
			this.agentRepository.getAgentsByPersonId(personId)
		]);
		if (personPropertyOwners.length > 0) roles.push('Owner');
		if (tenants.length > 0) roles.push('Tenant');
		if (households.length > 0) roles.push('Household');
		if (agents.length > 0) roles.push('Agent');
		return roles;
	}
}
