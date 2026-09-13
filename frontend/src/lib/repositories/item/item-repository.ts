import PocketBase from 'pocketbase';
import type { IItemModel } from '$lib/models/item-models';

export interface IItemRepository {
  getAllItems(): Promise<IItemModel[]>;
  getItemById(id: string): Promise<IItemModel>;
  createItem(name: string, description: string): Promise<IItemModel>;
  updateItem(id: string, name: string, description: string): Promise<void>;
  deleteItem(id: string): Promise<void>;
}

export class ItemRepository implements IItemRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getAllItems(): Promise<IItemModel[]> {
    return await this.backendClient.collection('items').getFullList<IItemModel>();
  }

  async getItemById(id: string): Promise<IItemModel> {
    // PocketBase v0.22: combining id= filter with a back-relation listRule
    // produces incorrect SQL. Fetch all accessible items and match client-side.
    const allItems = await this.backendClient.collection('items').getFullList<IItemModel>();
    const item = allItems.find((candidateItem) => candidateItem.id === id);
    if (!item) throw new Error(`Item not found: ${id}`);
    return item;
  }

  async createItem(name: string, description: string): Promise<IItemModel> {
    return await this.backendClient.collection('items').create<IItemModel>({ name, description });
  }

  async updateItem(id: string, name: string, description: string): Promise<void> {
    await this.backendClient.collection('items').update(id, { name, description });
  }

  async deleteItem(id: string): Promise<void> {
    await this.backendClient.collection('items').delete(id);
  }
}
