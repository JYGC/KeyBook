import PocketBase from 'pocketbase';
import type { IPersonItemModel } from '$lib/models/item-models';

export interface IPersonItemRepository {
  getByItemId(itemId: string): Promise<IPersonItemModel[]>;
  getByPersonId(personId: string): Promise<IPersonItemModel[]>;
  create(personId: string, itemId: string): Promise<IPersonItemModel>;
  delete(id: string): Promise<void>;
}

export class PersonItemRepository implements IPersonItemRepository {
  constructor(private readonly pb: PocketBase) {}

  async getByItemId(itemId: string): Promise<IPersonItemModel[]> {
    return await this.pb.collection('personItems').getFullList<IPersonItemModel>({
      filter: `item = "${itemId}"`,
    });
  }

  async getByPersonId(personId: string): Promise<IPersonItemModel[]> {
    return await this.pb.collection('personItems').getFullList<IPersonItemModel>({
      filter: `person = "${personId}"`,
    });
  }

  async create(personId: string, itemId: string): Promise<IPersonItemModel> {
    return await this.pb.collection('personItems').create<IPersonItemModel>({
      person: personId,
      item: itemId,
    });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('personItems').delete(id);
  }
}
