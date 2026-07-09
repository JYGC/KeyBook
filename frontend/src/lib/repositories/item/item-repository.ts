import PocketBase from 'pocketbase';
import type { IItemModel } from '$lib/models/item-models';

export interface IItemRepository {
  getAll(): Promise<IItemModel[]>;
  getById(id: string): Promise<IItemModel>;
  create(name: string, description: string): Promise<IItemModel>;
  update(id: string, name: string, description: string): Promise<void>;
  delete(id: string): Promise<void>;
}

export class ItemRepository implements IItemRepository {
  constructor(private readonly pb: PocketBase) {}

  async getAll(): Promise<IItemModel[]> {
    return await this.pb.collection('items').getFullList<IItemModel>();
  }

  async getById(id: string): Promise<IItemModel> {
    // PocketBase v0.22: combining id= filter with a back-relation listRule
    // produces incorrect SQL. Fetch all accessible items and match client-side.
    const all = await this.pb.collection('items').getFullList<IItemModel>();
    const item = all.find(i => i.id === id);
    if (!item) throw new Error(`Item not found: ${id}`);
    return item;
  }

  async create(name: string, description: string): Promise<IItemModel> {
    return await this.pb.collection('items').create<IItemModel>({ name, description });
  }

  async update(id: string, name: string, description: string): Promise<void> {
    await this.pb.collection('items').update(id, { name, description });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('items').delete(id);
  }
}
