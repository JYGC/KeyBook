import PocketBase from 'pocketbase';
import type { ITenantModel } from '$lib/models/property-models';

export interface ITenantRepository {
  getByPropertyId(propertyId: string): Promise<ITenantModel[]>;
  getByPersonId(personId: string): Promise<ITenantModel[]>;
  create(personId: string, propertyId: string): Promise<ITenantModel>;
  delete(id: string): Promise<void>;
}

export class TenantRepository implements ITenantRepository {
  constructor(private readonly pb: PocketBase) {}

  async getByPropertyId(propertyId: string): Promise<ITenantModel[]> {
    const all = await this.pb.collection('tenants').getFullList<ITenantModel>();
    return all.filter((r) => r.property === propertyId);
  }

  async getByPersonId(personId: string): Promise<ITenantModel[]> {
    const all = await this.pb.collection('tenants').getFullList<ITenantModel>();
    return all.filter((r) => r.person === personId);
  }

  async create(personId: string, propertyId: string): Promise<ITenantModel> {
    return await this.pb
      .collection('tenants')
      .create<ITenantModel>({ person: personId, property: propertyId });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('tenants').delete(id);
  }
}
