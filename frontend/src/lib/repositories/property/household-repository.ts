import PocketBase from 'pocketbase';
import type { IHouseholdModel } from '$lib/models/property-models';

export interface IHouseholdRepository {
  getByPropertyId(propertyId: string): Promise<IHouseholdModel[]>;
  getByPersonId(personId: string): Promise<IHouseholdModel[]>;
  create(personId: string, propertyId: string): Promise<IHouseholdModel>;
  delete(id: string): Promise<void>;
}

export class HouseholdRepository implements IHouseholdRepository {
  constructor(private readonly pb: PocketBase) {}

  async getByPropertyId(propertyId: string): Promise<IHouseholdModel[]> {
    const all = await this.pb.collection('households').getFullList<IHouseholdModel>();
    return all.filter((r) => r.property === propertyId);
  }

  async getByPersonId(personId: string): Promise<IHouseholdModel[]> {
    const all = await this.pb.collection('households').getFullList<IHouseholdModel>();
    return all.filter((r) => r.person === personId);
  }

  async create(personId: string, propertyId: string): Promise<IHouseholdModel> {
    return await this.pb
      .collection('households')
      .create<IHouseholdModel>({ person: personId, property: propertyId });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('households').delete(id);
  }
}
