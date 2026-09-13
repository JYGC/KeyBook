import PocketBase from 'pocketbase';
import type { IPropertyOwnerModel } from '$lib/models/property-models';

export interface IPropertyOwnerRepository {
  getPropertyOwnersByPropertyId(propertyId: string): Promise<IPropertyOwnerModel[]>;
  createPropertyOwner(propertyId: string): Promise<IPropertyOwnerModel>;
  deletePropertyOwner(id: string): Promise<void>;
}

export class PropertyOwnerRepository implements IPropertyOwnerRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getPropertyOwnersByPropertyId(propertyId: string): Promise<IPropertyOwnerModel[]> {
    const allPropertyOwners = await this.backendClient.collection('propertyOwners').getFullList<IPropertyOwnerModel>();
    return allPropertyOwners.filter((propertyOwner) => propertyOwner.property === propertyId);
  }

  async createPropertyOwner(propertyId: string): Promise<IPropertyOwnerModel> {
    return await this.backendClient
      .collection('propertyOwners')
      .create<IPropertyOwnerModel>({ property: propertyId });
  }

  async deletePropertyOwner(id: string): Promise<void> {
    await this.backendClient.collection('propertyOwners').delete(id);
  }
}
