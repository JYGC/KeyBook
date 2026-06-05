import PocketBase from 'pocketbase';
import type { IEntryDeviceModel } from '$lib/models/item-models';

export interface IEntryDeviceRepository {
  getById(id: string): Promise<IEntryDeviceModel>;
  getByItemId(itemId: string): Promise<IEntryDeviceModel | null>;
  create(itemId: string, deviceType: string, identifier: string, defunctReason: string): Promise<IEntryDeviceModel>;
  update(id: string, deviceType: string, identifier: string, defunctReason: string): Promise<void>;
  delete(id: string): Promise<void>;
}

export class EntryDeviceRepository implements IEntryDeviceRepository {
  constructor(private readonly pb: PocketBase) {}

  async getById(id: string): Promise<IEntryDeviceModel> {
    return await this.pb.collection('entryDevices').getOne<IEntryDeviceModel>(id);
  }

  async getByItemId(itemId: string): Promise<IEntryDeviceModel | null> {
    const results = await this.pb.collection('entryDevices').getFullList<IEntryDeviceModel>({
      filter: `item = "${itemId}"`,
    });
    return results.length > 0 ? results[0] : null;
  }

  async create(
    itemId: string,
    deviceType: string,
    identifier: string,
    defunctReason: string,
  ): Promise<IEntryDeviceModel> {
    return await this.pb.collection('entryDevices').create<IEntryDeviceModel>({
      item: itemId,
      deviceType,
      identifier,
      defunctReason,
    });
  }

  async update(id: string, deviceType: string, identifier: string, defunctReason: string): Promise<void> {
    await this.pb.collection('entryDevices').update(id, { deviceType, identifier, defunctReason });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('entryDevices').delete(id);
  }
}
