import PocketBase from 'pocketbase';
import type { ITenantModel } from '$lib/models/property-models';

export interface ITenantRepository {
  getTenantsByPropertyId(propertyId: string): Promise<ITenantModel[]>;
  getTenantsByPersonId(personId: string): Promise<ITenantModel[]>;
  createTenant(personId: string, propertyId: string): Promise<ITenantModel>;
  deleteTenant(id: string): Promise<void>;
}

export class TenantRepository implements ITenantRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getTenantsByPropertyId(propertyId: string): Promise<ITenantModel[]> {
    const allTenants = await this.backendClient.collection('tenants').getFullList<ITenantModel>();
    return allTenants.filter((tenant) => tenant.property === propertyId);
  }

  async getTenantsByPersonId(personId: string): Promise<ITenantModel[]> {
    const allTenants = await this.backendClient.collection('tenants').getFullList<ITenantModel>();
    return allTenants.filter((tenant) => tenant.person === personId);
  }

  async createTenant(personId: string, propertyId: string): Promise<ITenantModel> {
    return await this.backendClient
      .collection('tenants')
      .create<ITenantModel>({ person: personId, property: propertyId });
  }

  async deleteTenant(id: string): Promise<void> {
    await this.backendClient.collection('tenants').delete(id);
  }
}
