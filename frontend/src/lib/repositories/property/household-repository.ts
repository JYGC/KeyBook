import PocketBase from 'pocketbase';
import type { IHouseholdModel } from '$lib/models/property-models';

export interface IHouseholdRepository {
  getHouseholdsByPropertyId(propertyId: string): Promise<IHouseholdModel[]>;
  getHouseholdsByPersonId(personId: string): Promise<IHouseholdModel[]>;
  createHousehold(personId: string, propertyId: string): Promise<IHouseholdModel>;
  deleteHousehold(id: string): Promise<void>;
}

export class HouseholdRepository implements IHouseholdRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getHouseholdsByPropertyId(propertyId: string): Promise<IHouseholdModel[]> {
    const allHouseholds = await this.backendClient.collection('households').getFullList<IHouseholdModel>();
    return allHouseholds.filter((household) => household.property === propertyId);
  }

  async getHouseholdsByPersonId(personId: string): Promise<IHouseholdModel[]> {
    const allHouseholds = await this.backendClient.collection('households').getFullList<IHouseholdModel>();
    return allHouseholds.filter((household) => household.person === personId);
  }

  async createHousehold(personId: string, propertyId: string): Promise<IHouseholdModel> {
    return await this.backendClient
      .collection('households')
      .create<IHouseholdModel>({ person: personId, property: propertyId });
  }

  async deleteHousehold(id: string): Promise<void> {
    await this.backendClient.collection('households').delete(id);
  }
}
