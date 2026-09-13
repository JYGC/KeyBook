import PocketBase from 'pocketbase';
import type { IPersonItemModel } from '$lib/models/item-models';

export interface IPersonItemRepository {
  getPersonItemsByItemId(itemId: string): Promise<IPersonItemModel[]>;
  getPersonItemsByPersonId(personId: string): Promise<IPersonItemModel[]>;
  createPersonItem(personId: string, itemId: string): Promise<IPersonItemModel>;
  deletePersonItem(id: string): Promise<void>;
}

export class PersonItemRepository implements IPersonItemRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getPersonItemsByItemId(itemId: string): Promise<IPersonItemModel[]> {
    return await this.backendClient.collection('personItems').getFullList<IPersonItemModel>({
      filter: `item = "${itemId}"`,
    });
  }

  async getPersonItemsByPersonId(personId: string): Promise<IPersonItemModel[]> {
    return await this.backendClient.collection('personItems').getFullList<IPersonItemModel>({
      filter: `person = "${personId}"`,
    });
  }

  async createPersonItem(personId: string, itemId: string): Promise<IPersonItemModel> {
    return await this.backendClient.collection('personItems').create<IPersonItemModel>({
      person: personId,
      item: itemId,
    });
  }

  async deletePersonItem(id: string): Promise<void> {
    await this.backendClient.collection('personItems').delete(id);
  }
}
