import PocketBase from 'pocketbase';
import type { ICobrandPropertyManagerModel } from '$lib/models/cobrand-models';

export interface ICobrandPropertyManagerRepository {
  getByCobrandId(cobrandId: string): Promise<ICobrandPropertyManagerModel[]>;
  getByPropertyId(propertyId: string): Promise<ICobrandPropertyManagerModel[]>;
  create(cobrandId: string, propertyId: string): Promise<ICobrandPropertyManagerModel>;
  delete(id: string): Promise<void>;
}

export class CobrandPropertyManagerRepository implements ICobrandPropertyManagerRepository {
  constructor(private readonly pb: PocketBase) {}

  async getByCobrandId(cobrandId: string): Promise<ICobrandPropertyManagerModel[]> {
    const all = await this.pb
      .collection('cobrandPropertyManagers')
      .getFullList<ICobrandPropertyManagerModel>();
    return all.filter((m) => m.cobrand === cobrandId);
  }

  async getByPropertyId(propertyId: string): Promise<ICobrandPropertyManagerModel[]> {
    const all = await this.pb
      .collection('cobrandPropertyManagers')
      .getFullList<ICobrandPropertyManagerModel>();
    return all.filter((m) => m.property === propertyId);
  }

  async create(cobrandId: string, propertyId: string): Promise<ICobrandPropertyManagerModel> {
    return await this.pb
      .collection('cobrandPropertyManagers')
      .create<ICobrandPropertyManagerModel>({
        cobrand: cobrandId,
        property: propertyId,
      });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('cobrandPropertyManagers').delete(id);
  }
}
