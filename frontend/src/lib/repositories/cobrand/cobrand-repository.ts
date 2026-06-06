import PocketBase from 'pocketbase';
import type { ICobrandModel } from '$lib/models/cobrand-models';

export interface ICobrandRepository {
  getAll(): Promise<ICobrandModel[]>;
  getById(id: string): Promise<ICobrandModel>;
  create(name: string): Promise<ICobrandModel>;
  update(id: string, name: string): Promise<void>;
  delete(id: string): Promise<void>;
}

export class CobrandRepository implements ICobrandRepository {
  constructor(private readonly pb: PocketBase) {}

  async getAll(): Promise<ICobrandModel[]> {
    return await this.pb.collection('cobrands').getFullList<ICobrandModel>();
  }

  async getById(id: string): Promise<ICobrandModel> {
    // PocketBase v0.22: combining id= filter with a back-relation listRule
    // produces incorrect SQL. Fetch all accessible records and match client-side.
    const all = await this.pb.collection('cobrands').getFullList<ICobrandModel>();
    const cobrand = all.find((c) => c.id === id);
    if (!cobrand) throw new Error(`Cobrand not found: ${id}`);
    return cobrand;
  }

  async create(name: string): Promise<ICobrandModel> {
    return await this.pb.collection('cobrands').create<ICobrandModel>({ name });
  }

  async update(id: string, name: string): Promise<void> {
    await this.pb.collection('cobrands').update(id, { name });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('cobrands').delete(id);
  }
}
