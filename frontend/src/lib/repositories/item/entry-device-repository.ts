import PocketBase from 'pocketbase';
import type { IEntryDeviceModel } from '$lib/models/item-models';

export interface IEntryDeviceRepository {
  getEntryDeviceById(id: string): Promise<IEntryDeviceModel>;
  getEntryDeviceByItemId(itemId: string): Promise<IEntryDeviceModel | null>;
  createEntryDevice(itemId: string, deviceType: string, identifier: string, defunctReason: string): Promise<IEntryDeviceModel>;
  updateEntryDevice(id: string, deviceType: string, identifier: string, defunctReason: string): Promise<void>;
  deleteEntryDevice(id: string): Promise<void>;
}

export class EntryDeviceRepository implements IEntryDeviceRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getEntryDeviceById(id: string): Promise<IEntryDeviceModel> {
    return await this.backendClient.collection('entryDevices').getOne<IEntryDeviceModel>(id);
  }

  async getEntryDeviceByItemId(itemId: string): Promise<IEntryDeviceModel | null> {
    // PocketBase v0.22: combining an explicit item= filter with the listRule (which also
    // traverses item) produces incorrect SQL. Fetch all and match client-side.
    const allEntryDevices = await this.backendClient.collection('entryDevices').getFullList<IEntryDeviceModel>();
    return allEntryDevices.find((entryDevice) => entryDevice.item === itemId) ?? null;
  }

  async createEntryDevice(
    itemId: string,
    deviceType: string,
    identifier: string,
    defunctReason: string,
  ): Promise<IEntryDeviceModel> {
    return await this.backendClient.collection('entryDevices').create<IEntryDeviceModel>({
      item: itemId,
      deviceType,
      identifier,
      defunctReason,
    });
  }

  async updateEntryDevice(id: string, deviceType: string, identifier: string, defunctReason: string): Promise<void> {
    await this.backendClient.collection('entryDevices').update(id, { deviceType, identifier, defunctReason });
  }

  async deleteEntryDevice(id: string): Promise<void> {
    await this.backendClient.collection('entryDevices').delete(id);
  }
}
