import PocketBase from 'pocketbase';
import type { ICobrandPropertyOwnerModel } from '$lib/models/cobrand-models';

export interface ICobrandPropertyOwnerRepository {
  getByCobrandId(cobrandId: string): Promise<ICobrandPropertyOwnerModel[]>;
  getByPropertyOwnerId(propertyOwnerId: string): Promise<ICobrandPropertyOwnerModel[]>;
  create(cobrandId: string, propertyOwnerId: string): Promise<ICobrandPropertyOwnerModel>;
  delete(id: string): Promise<void>;
}

export class CobrandPropertyOwnerRepository implements ICobrandPropertyOwnerRepository {
  constructor(private readonly pb: PocketBase) {}

  async getByCobrandId(cobrandId: string): Promise<ICobrandPropertyOwnerModel[]> {
    const all = await this.pb
      .collection('cobrandPropertyOwners')
      .getFullList<ICobrandPropertyOwnerModel>();
    return all.filter((o) => o.cobrand === cobrandId);
  }

  async getByPropertyOwnerId(propertyOwnerId: string): Promise<ICobrandPropertyOwnerModel[]> {
    const all = await this.pb
      .collection('cobrandPropertyOwners')
      .getFullList<ICobrandPropertyOwnerModel>();
    return all.filter((o) => o.propertyOwner === propertyOwnerId);
  }

  async create(cobrandId: string, propertyOwnerId: string): Promise<ICobrandPropertyOwnerModel> {
    return await this.pb
      .collection('cobrandPropertyOwners')
      .create<ICobrandPropertyOwnerModel>({
        cobrand: cobrandId,
        propertyOwner: propertyOwnerId,
      });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('cobrandPropertyOwners').delete(id);
  }
}
