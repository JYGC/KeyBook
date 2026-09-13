import PocketBase from 'pocketbase';
import type { IPropertyItemModel } from '$lib/models/item-models';

export interface IPropertyItemRepository {
  getPropertyItemsByItemId(itemId: string): Promise<IPropertyItemModel[]>;
  getPropertyItemsByPropertyId(propertyId: string): Promise<IPropertyItemModel[]>;
  createPropertyItem(itemId: string, propertyId: string): Promise<IPropertyItemModel>;
  deletePropertyItem(id: string): Promise<void>;
}

export class PropertyItemRepository implements IPropertyItemRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getPropertyItemsByItemId(itemId: string): Promise<IPropertyItemModel[]> {
    return await this.backendClient.collection('propertyItems').getFullList<IPropertyItemModel>({
      filter: `item = "${itemId}"`,
    });
  }

  async getPropertyItemsByPropertyId(propertyId: string): Promise<IPropertyItemModel[]> {
    return await this.backendClient.collection('propertyItems').getFullList<IPropertyItemModel>({
      filter: `property = "${propertyId}"`,
    });
  }

  async createPropertyItem(itemId: string, propertyId: string): Promise<IPropertyItemModel> {
    return await this.backendClient.collection('propertyItems').create<IPropertyItemModel>({
      item: itemId,
      property: propertyId,
    });
  }

  async deletePropertyItem(id: string): Promise<void> {
    await this.backendClient.collection('propertyItems').delete(id);
  }
}
