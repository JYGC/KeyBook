import PocketBase from 'pocketbase';
import type { ICobrandPropertyOwnerModel } from '$lib/models/cobrand-models';

export interface ICobrandPropertyOwnerRepository {
  getCobrandPropertyOwnersByCobrandId(cobrandId: string): Promise<ICobrandPropertyOwnerModel[]>;
  getCobrandPropertyOwnersByPropertyOwnerId(propertyOwnerId: string): Promise<ICobrandPropertyOwnerModel[]>;
  createCobrandPropertyOwner(cobrandId: string, propertyOwnerId: string): Promise<ICobrandPropertyOwnerModel>;
  deleteCobrandPropertyOwner(id: string): Promise<void>;
}

export class CobrandPropertyOwnerRepository implements ICobrandPropertyOwnerRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getCobrandPropertyOwnersByCobrandId(cobrandId: string): Promise<ICobrandPropertyOwnerModel[]> {
    const allCobrandPropertyOwners = await this.backendClient
      .collection('cobrandPropertyOwners')
      .getFullList<ICobrandPropertyOwnerModel>();
    return allCobrandPropertyOwners.filter((cobrandPropertyOwner) => cobrandPropertyOwner.cobrand === cobrandId);
  }

  async getCobrandPropertyOwnersByPropertyOwnerId(propertyOwnerId: string): Promise<ICobrandPropertyOwnerModel[]> {
    const allCobrandPropertyOwners = await this.backendClient
      .collection('cobrandPropertyOwners')
      .getFullList<ICobrandPropertyOwnerModel>();
    return allCobrandPropertyOwners.filter((cobrandPropertyOwner) => cobrandPropertyOwner.propertyOwner === propertyOwnerId);
  }

  async createCobrandPropertyOwner(cobrandId: string, propertyOwnerId: string): Promise<ICobrandPropertyOwnerModel> {
    return await this.backendClient
      .collection('cobrandPropertyOwners')
      .create<ICobrandPropertyOwnerModel>({
        cobrand: cobrandId,
        propertyOwner: propertyOwnerId,
      });
  }

  async deleteCobrandPropertyOwner(id: string): Promise<void> {
    await this.backendClient.collection('cobrandPropertyOwners').delete(id);
  }
}
