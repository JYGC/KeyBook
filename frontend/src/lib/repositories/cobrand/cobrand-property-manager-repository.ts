import PocketBase from 'pocketbase';
import type { ICobrandPropertyManagerModel } from '$lib/models/cobrand-models';

export interface ICobrandPropertyManagerRepository {
  getCobrandPropertyManagersByCobrandId(cobrandId: string): Promise<ICobrandPropertyManagerModel[]>;
  getCobrandPropertyManagersByPropertyId(propertyId: string): Promise<ICobrandPropertyManagerModel[]>;
  createCobrandPropertyManager(cobrandId: string, propertyId: string): Promise<ICobrandPropertyManagerModel>;
  deleteCobrandPropertyManager(id: string): Promise<void>;
}

export class CobrandPropertyManagerRepository implements ICobrandPropertyManagerRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getCobrandPropertyManagersByCobrandId(cobrandId: string): Promise<ICobrandPropertyManagerModel[]> {
    const allCobrandPropertyManagers = await this.backendClient
      .collection('cobrandPropertyManagers')
      .getFullList<ICobrandPropertyManagerModel>();
    return allCobrandPropertyManagers.filter((cobrandPropertyManager) => cobrandPropertyManager.cobrand === cobrandId);
  }

  async getCobrandPropertyManagersByPropertyId(propertyId: string): Promise<ICobrandPropertyManagerModel[]> {
    const allCobrandPropertyManagers = await this.backendClient
      .collection('cobrandPropertyManagers')
      .getFullList<ICobrandPropertyManagerModel>();
    return allCobrandPropertyManagers.filter((cobrandPropertyManager) => cobrandPropertyManager.property === propertyId);
  }

  async createCobrandPropertyManager(cobrandId: string, propertyId: string): Promise<ICobrandPropertyManagerModel> {
    return await this.backendClient
      .collection('cobrandPropertyManagers')
      .create<ICobrandPropertyManagerModel>({
        cobrand: cobrandId,
        property: propertyId,
      });
  }

  async deleteCobrandPropertyManager(id: string): Promise<void> {
    await this.backendClient.collection('cobrandPropertyManagers').delete(id);
  }
}
