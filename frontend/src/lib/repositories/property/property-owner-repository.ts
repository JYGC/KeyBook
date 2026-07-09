import PocketBase from 'pocketbase';
import type { IPropertyOwnerModel } from '$lib/models/property-models';

export interface IPropertyOwnerRepository {
  getByPropertyId(propertyId: string): Promise<IPropertyOwnerModel[]>;
  create(propertyId: string): Promise<IPropertyOwnerModel>;
  delete(id: string): Promise<void>;
}

export class PropertyOwnerRepository implements IPropertyOwnerRepository {
  constructor(private readonly pb: PocketBase) {}

  async getByPropertyId(propertyId: string): Promise<IPropertyOwnerModel[]> {
    const all = await this.pb.collection('propertyOwners').getFullList<IPropertyOwnerModel>();
    return all.filter((r) => r.property === propertyId);
  }

  async create(propertyId: string): Promise<IPropertyOwnerModel> {
    return await this.pb
      .collection('propertyOwners')
      .create<IPropertyOwnerModel>({ property: propertyId });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('propertyOwners').delete(id);
  }
}
