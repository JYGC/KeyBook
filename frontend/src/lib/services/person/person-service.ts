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
  createPerson(name: string, dob: string): Promise<IPersonModel>;
  updatePerson(id: string, name: string, dob: string): Promise<IPersonModel>;
  deletePerson(id: string): Promise<void>;
  getRolesForPerson(personId: string): Promise<string[]>;
}

export class PersonService implements IPersonService {
  constructor(
    private readonly personRepo: IPersonRepository,
    private readonly ppoRepo: IPersonPropertyOwnerRepository,
    private readonly tenantRepo: ITenantRepository,
    private readonly householdRepo: IHouseholdRepository,
    private readonly agentRepo: IAgentRepository,
  ) {}

  validatePersonName(name: string): void {
    if (!name || name.trim() === '') throw new Error('person name is required');
  }

  validatePersonDOB(dob: string): void {
    if (!dob || dob.trim() === '') throw new Error('date of birth is required');
  }

  async getAllPersons(): Promise<IPersonModel[]> {
    return await this.personRepo.getAll();
  }

  async getPersonById(id: string): Promise<IPersonModel> {
    return await this.personRepo.getById(id);
  }

  async getPersonByUserId(userId: string): Promise<IPersonModel | null> {
    return await this.personRepo.getByUserId(userId);
  }

  async createPerson(name: string, dob: string): Promise<IPersonModel> {
    this.validatePersonName(name);
    this.validatePersonDOB(dob);
    return await this.personRepo.create(name, dob);
  }

  async updatePerson(id: string, name: string, dob: string): Promise<IPersonModel> {
    this.validatePersonName(name);
    this.validatePersonDOB(dob);
    return await this.personRepo.update(id, name, dob);
  }

  async deletePerson(id: string): Promise<void> {
    await this.personRepo.delete(id);
  }

  async getRolesForPerson(personId: string): Promise<string[]> {
    const roles: string[] = [];
    const [ppos, tenants, households, agents] = await Promise.all([
      this.ppoRepo.getByPersonId(personId),
      this.tenantRepo.getByPersonId(personId),
      this.householdRepo.getByPersonId(personId),
      this.agentRepo.getByPersonId(personId),
    ]);
    if (ppos.length > 0) roles.push('Owner');
    if (tenants.length > 0) roles.push('Tenant');
    if (households.length > 0) roles.push('Household');
    if (agents.length > 0) roles.push('Agent');
    return roles;
  }
}
