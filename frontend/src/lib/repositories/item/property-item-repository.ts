import PocketBase from 'pocketbase';
import type { IPropertyItemModel } from '$lib/models/item-models';

export interface IPropertyItemRepository {
  getByItemId(itemId: string): Promise<IPropertyItemModel[]>;
  getByPropertyId(propertyId: string): Promise<IPropertyItemModel[]>;
  create(itemId: string, propertyId: string): Promise<IPropertyItemModel>;
  delete(id: string): Promise<void>;
}

export class PropertyItemRepository implements IPropertyItemRepository {
  constructor(private readonly pb: PocketBase) {}

  async getByItemId(itemId: string): Promise<IPropertyItemModel[]> {
    return await this.pb.collection('propertyItems').getFullList<IPropertyItemModel>({
      filter: `item = "${itemId}"`,
    });
  }

  async getByPropertyId(propertyId: string): Promise<IPropertyItemModel[]> {
    return await this.pb.collection('propertyItems').getFullList<IPropertyItemModel>({
      filter: `property = "${propertyId}"`,
    });
  }

  async create(itemId: string, propertyId: string): Promise<IPropertyItemModel> {
    return await this.pb.collection('propertyItems').create<IPropertyItemModel>({
      item: itemId,
      property: propertyId,
    });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('propertyItems').delete(id);
  }
}
