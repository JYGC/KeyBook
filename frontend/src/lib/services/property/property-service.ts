import type { IPropertyRepository } from '$lib/repositories/property/property-repository';
import type { IPropertyOwnerRepository } from '$lib/repositories/property/property-owner-repository';
import type { IPersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
import type { ITenantRepository } from '$lib/repositories/property/tenant-repository';
import type { IHouseholdRepository } from '$lib/repositories/property/household-repository';
import type { IPropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';
import type { IPropertyModel, IPersonPropertyOwnerModel, ITenantModel, IHouseholdModel } from '$lib/models/property-models';
import type { IPropertyAgentModel } from '$lib/models/agent-models';

export interface IPropertyService {
  validatePropertyAddress(address: string): void;
  getAllProperties(): Promise<IPropertyModel[]>;
  getPropertyById(id: string): Promise<IPropertyModel>;
  createPropertyWithOwner(address: string, personId: string): Promise<IPropertyModel>;
  updateProperty(id: string, address: string): Promise<IPropertyModel>;
  deleteProperty(id: string): Promise<void>;
  getPersonOwnersForProperty(propertyId: string): Promise<IPersonPropertyOwnerModel[]>;
  getTenantsForProperty(propertyId: string): Promise<ITenantModel[]>;
  getHouseholdMembersForProperty(propertyId: string): Promise<IHouseholdModel[]>;
  getPropertyAgentsForProperty(propertyId: string): Promise<IPropertyAgentModel[]>;
  addTenant(personId: string, propertyId: string): Promise<ITenantModel>;
  removeTenant(id: string): Promise<void>;
  addHouseholdMember(personId: string, propertyId: string): Promise<IHouseholdModel>;
  removeHouseholdMember(id: string): Promise<void>;
}

export class PropertyService implements IPropertyService {
  constructor(
    private readonly propertyRepo: IPropertyRepository,
    private readonly propertyOwnerRepo: IPropertyOwnerRepository,
    private readonly ppoRepo: IPersonPropertyOwnerRepository,
    private readonly tenantRepo: ITenantRepository,
    private readonly householdRepo: IHouseholdRepository,
    private readonly propertyAgentRepo: IPropertyAgentRepository,
  ) {}

  validatePropertyAddress(address: string): void {
    if (!address || address.trim() === '') throw new Error('property address is required');
  }

  async getAllProperties(): Promise<IPropertyModel[]> {
    return await this.propertyRepo.getAll();
  }

  async getPropertyById(id: string): Promise<IPropertyModel> {
    return await this.propertyRepo.getById(id);
  }

  async createPropertyWithOwner(address: string, personId: string): Promise<IPropertyModel> {
    this.validatePropertyAddress(address);
    const property = await this.propertyRepo.create(address);
    const propertyOwner = await this.propertyOwnerRepo.create(property.id);
    await this.ppoRepo.create(personId, propertyOwner.id);
    return property;
  }

  async updateProperty(id: string, address: string): Promise<IPropertyModel> {
    this.validatePropertyAddress(address);
    return await this.propertyRepo.update(id, address);
  }

  async deleteProperty(id: string): Promise<void> {
    await this.propertyRepo.delete(id);
  }

  async getPersonOwnersForProperty(propertyId: string): Promise<IPersonPropertyOwnerModel[]> {
    const propertyOwners = await this.propertyOwnerRepo.getByPropertyId(propertyId);
    const results: IPersonPropertyOwnerModel[] = [];
    for (const po of propertyOwners) {
      const ppos = await this.ppoRepo.getByPropertyOwnerId(po.id);
      results.push(...ppos);
    }
    return results;
  }

  async getTenantsForProperty(propertyId: string): Promise<ITenantModel[]> {
    return await this.tenantRepo.getByPropertyId(propertyId);
  }

  async getHouseholdMembersForProperty(propertyId: string): Promise<IHouseholdModel[]> {
    return await this.householdRepo.getByPropertyId(propertyId);
  }

  async getPropertyAgentsForProperty(propertyId: string): Promise<IPropertyAgentModel[]> {
    return await this.propertyAgentRepo.getByPropertyId(propertyId);
  }

  async addTenant(personId: string, propertyId: string): Promise<ITenantModel> {
    return await this.tenantRepo.create(personId, propertyId);
  }

  async removeTenant(id: string): Promise<void> {
    await this.tenantRepo.delete(id);
  }

  async addHouseholdMember(personId: string, propertyId: string): Promise<IHouseholdModel> {
    return await this.householdRepo.create(personId, propertyId);
  }

  async removeHouseholdMember(id: string): Promise<void> {
    await this.householdRepo.delete(id);
  }
}
