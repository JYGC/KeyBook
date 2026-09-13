import PocketBase from 'pocketbase';
import type { IPersonPropertyOwnerModel } from '$lib/models/property-models';

export interface IPersonPropertyOwnerRepository {
  getPersonPropertyOwnersByPersonId(personId: string): Promise<IPersonPropertyOwnerModel[]>;
  getPersonPropertyOwnersByPropertyOwnerId(propertyOwnerId: string): Promise<IPersonPropertyOwnerModel[]>;
  createPersonPropertyOwner(personId: string, propertyOwnerId: string): Promise<IPersonPropertyOwnerModel>;
  deletePersonPropertyOwner(id: string): Promise<void>;
}

export class PersonPropertyOwnerRepository implements IPersonPropertyOwnerRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getPersonPropertyOwnersByPersonId(personId: string): Promise<IPersonPropertyOwnerModel[]> {
    const allPersonPropertyOwners = await this.backendClient
      .collection('personPropertyOwners')
      .getFullList<IPersonPropertyOwnerModel>();
    return allPersonPropertyOwners.filter((personPropertyOwner) => personPropertyOwner.person === personId);
  }

  async getPersonPropertyOwnersByPropertyOwnerId(propertyOwnerId: string): Promise<IPersonPropertyOwnerModel[]> {
    const allPersonPropertyOwners = await this.backendClient
      .collection('personPropertyOwners')
      .getFullList<IPersonPropertyOwnerModel>();
    return allPersonPropertyOwners.filter((personPropertyOwner) => personPropertyOwner.propertyOwner === propertyOwnerId);
  }

  async createPersonPropertyOwner(personId: string, propertyOwnerId: string): Promise<IPersonPropertyOwnerModel> {
    return await this.backendClient
      .collection('personPropertyOwners')
      .create<IPersonPropertyOwnerModel>({ person: personId, propertyOwner: propertyOwnerId });
  }

  async deletePersonPropertyOwner(id: string): Promise<void> {
    await this.backendClient.collection('personPropertyOwners').delete(id);
  }
}
