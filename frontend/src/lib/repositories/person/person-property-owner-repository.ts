import PocketBase from 'pocketbase';
import type { IPersonPropertyOwnerModel } from '$lib/models/property-models';

export interface IPersonPropertyOwnerRepository {
  getByPersonId(personId: string): Promise<IPersonPropertyOwnerModel[]>;
  getByPropertyOwnerId(propertyOwnerId: string): Promise<IPersonPropertyOwnerModel[]>;
  create(personId: string, propertyOwnerId: string): Promise<IPersonPropertyOwnerModel>;
  delete(id: string): Promise<void>;
}

export class PersonPropertyOwnerRepository implements IPersonPropertyOwnerRepository {
  constructor(private readonly pb: PocketBase) {}

  async getByPersonId(personId: string): Promise<IPersonPropertyOwnerModel[]> {
    const all = await this.pb
      .collection('personPropertyOwners')
      .getFullList<IPersonPropertyOwnerModel>();
    return all.filter((r) => r.person === personId);
  }

  async getByPropertyOwnerId(propertyOwnerId: string): Promise<IPersonPropertyOwnerModel[]> {
    const all = await this.pb
      .collection('personPropertyOwners')
      .getFullList<IPersonPropertyOwnerModel>();
    return all.filter((r) => r.propertyOwner === propertyOwnerId);
  }

  async create(personId: string, propertyOwnerId: string): Promise<IPersonPropertyOwnerModel> {
    return await this.pb
      .collection('personPropertyOwners')
      .create<IPersonPropertyOwnerModel>({ person: personId, propertyOwner: propertyOwnerId });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('personPropertyOwners').delete(id);
  }
}
