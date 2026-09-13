import PocketBase from 'pocketbase';
import type { ICobrandModel } from '$lib/models/cobrand-models';

export interface ICobrandRepository {
  getAllCobrands(): Promise<ICobrandModel[]>;
  getCobrandById(id: string): Promise<ICobrandModel>;
  createCobrand(name: string): Promise<ICobrandModel>;
  updateCobrand(id: string, name: string): Promise<void>;
  deleteCobrand(id: string): Promise<void>;
}

export class CobrandRepository implements ICobrandRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getAllCobrands(): Promise<ICobrandModel[]> {
    return await this.backendClient.collection('cobrands').getFullList<ICobrandModel>();
  }

  async getCobrandById(id: string): Promise<ICobrandModel> {
    // PocketBase v0.22: combining id= filter with a back-relation listRule
    // produces incorrect SQL. Fetch all accessible records and match client-side.
    const allCobrands = await this.backendClient.collection('cobrands').getFullList<ICobrandModel>();
    const cobrand = allCobrands.find((candidateCobrand) => candidateCobrand.id === id);
    if (!cobrand) throw new Error(`Cobrand not found: ${id}`);
    return cobrand;
  }

  async createCobrand(name: string): Promise<ICobrandModel> {
    return await this.backendClient.collection('cobrands').create<ICobrandModel>({ name });
  }

  async updateCobrand(id: string, name: string): Promise<void> {
    await this.backendClient.collection('cobrands').update(id, { name });
  }

  async deleteCobrand(id: string): Promise<void> {
    await this.backendClient.collection('cobrands').delete(id);
  }
}
