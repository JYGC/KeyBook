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
    private readonly propertyRepository: IPropertyRepository,
    private readonly propertyOwnerRepository: IPropertyOwnerRepository,
    private readonly personPropertyOwnerRepository: IPersonPropertyOwnerRepository,
    private readonly tenantRepository: ITenantRepository,
    private readonly householdRepository: IHouseholdRepository,
    private readonly propertyAgentRepository: IPropertyAgentRepository,
  ) {}

  validatePropertyAddress(address: string): void {
    if (!address || address.trim() === '') throw new Error('property address is required');
  }

  async getAllProperties(): Promise<IPropertyModel[]> {
    return await this.propertyRepository.getAllProperties();
  }

  async getPropertyById(id: string): Promise<IPropertyModel> {
    return await this.propertyRepository.getPropertyById(id);
  }

  async createPropertyWithOwner(address: string, personId: string): Promise<IPropertyModel> {
    this.validatePropertyAddress(address);
    const property = await this.propertyRepository.createProperty(address);
    const propertyOwner = await this.propertyOwnerRepository.createPropertyOwner(property.id);
    await this.personPropertyOwnerRepository.createPersonPropertyOwner(personId, propertyOwner.id);
    return property;
  }

  async updateProperty(id: string, address: string): Promise<IPropertyModel> {
    this.validatePropertyAddress(address);
    return await this.propertyRepository.updateProperty(id, address);
  }

  async deleteProperty(id: string): Promise<void> {
    await this.propertyRepository.deleteProperty(id);
  }

  async getPersonOwnersForProperty(propertyId: string): Promise<IPersonPropertyOwnerModel[]> {
    const propertyOwners = await this.propertyOwnerRepository.getPropertyOwnersByPropertyId(propertyId);
    const personOwners: IPersonPropertyOwnerModel[] = [];
    for (const propertyOwner of propertyOwners) {
      const personOwnersOfThisRecord =
        await this.personPropertyOwnerRepository.getPersonPropertyOwnersByPropertyOwnerId(
          propertyOwner.id,
        );
      personOwners.push(...personOwnersOfThisRecord);
    }
    return personOwners;
  }

  async getTenantsForProperty(propertyId: string): Promise<ITenantModel[]> {
    return await this.tenantRepository.getTenantsByPropertyId(propertyId);
  }

  async getHouseholdMembersForProperty(propertyId: string): Promise<IHouseholdModel[]> {
    return await this.householdRepository.getHouseholdsByPropertyId(propertyId);
  }

  async getPropertyAgentsForProperty(propertyId: string): Promise<IPropertyAgentModel[]> {
    return await this.propertyAgentRepository.getPropertyAgentsByPropertyId(propertyId);
  }

  async addTenant(personId: string, propertyId: string): Promise<ITenantModel> {
    return await this.tenantRepository.createTenant(personId, propertyId);
  }

  async removeTenant(id: string): Promise<void> {
    await this.tenantRepository.deleteTenant(id);
  }

  async addHouseholdMember(personId: string, propertyId: string): Promise<IHouseholdModel> {
    return await this.householdRepository.createHousehold(personId, propertyId);
  }

  async removeHouseholdMember(id: string): Promise<void> {
    await this.householdRepository.deleteHousehold(id);
  }
}
