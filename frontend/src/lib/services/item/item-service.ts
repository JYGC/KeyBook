import type { IItemRepository } from '$lib/repositories/item/item-repository';
import type { IEntryDeviceRepository } from '$lib/repositories/item/entry-device-repository';
import type { IItemModel, IEntryDeviceModel } from '$lib/models/item-models';

export interface IItemService {
  validateItemName(name: string): void;
  validateEntryDeviceTransition(currentDefunctReason: string, newDefunctReason: string): void;
  getAllItems(): Promise<IItemModel[]>;
  getItemWithEntryDevice(itemId: string): Promise<{ item: IItemModel; entryDevice: IEntryDeviceModel | null }>;
  createItem(name: string, description: string): Promise<IItemModel>;
  updateItem(id: string, name: string, description: string): Promise<void>;
  deleteItem(id: string): Promise<void>;
  createEntryDevice(itemId: string, deviceType: string, identifier: string, defunctReason: string): Promise<IEntryDeviceModel>;
  updateEntryDevice(id: string, deviceType: string, identifier: string, newDefunctReason: string): Promise<void>;
  deleteEntryDevice(id: string): Promise<void>;
}

export class ItemService implements IItemService {
  constructor(
    private readonly itemRepo: IItemRepository,
    private readonly entryDeviceRepo: IEntryDeviceRepository,
  ) {}

  validateItemName(name: string): void {
    if (!name || name.trim() === '') {
      throw new Error('item name is required');
    }
  }

  validateEntryDeviceTransition(currentDefunctReason: string, newDefunctReason: string): void {
    const isActive = (r: string) => r === '' || r === 'None';
    if (!isActive(currentDefunctReason) && isActive(newDefunctReason)) {
      throw new Error('cannot reactivate a defunct entry device');
    }
  }

  async getAllItems(): Promise<IItemModel[]> {
    return await this.itemRepo.getAll();
  }

  async getItemWithEntryDevice(
    itemId: string,
  ): Promise<{ item: IItemModel; entryDevice: IEntryDeviceModel | null }> {
    const [item, entryDevice] = await Promise.all([
      this.itemRepo.getById(itemId),
      this.entryDeviceRepo.getByItemId(itemId),
    ]);
    return { item, entryDevice };
  }

  async createItem(name: string, description: string): Promise<IItemModel> {
    this.validateItemName(name);
    return await this.itemRepo.create(name, description);
  }

  async updateItem(id: string, name: string, description: string): Promise<void> {
    this.validateItemName(name);
    await this.itemRepo.update(id, name, description);
  }

  async deleteItem(id: string): Promise<void> {
    const entryDevice = await this.entryDeviceRepo.getByItemId(id);
    if (entryDevice !== null) {
      await this.entryDeviceRepo.delete(entryDevice.id);
    }
    await this.itemRepo.delete(id);
  }

  async createEntryDevice(
    itemId: string,
    deviceType: string,
    identifier: string,
    defunctReason: string,
  ): Promise<IEntryDeviceModel> {
    return await this.entryDeviceRepo.create(itemId, deviceType, identifier, defunctReason);
  }

  async updateEntryDevice(
    id: string,
    deviceType: string,
    identifier: string,
    newDefunctReason: string,
  ): Promise<void> {
    const current = await this.entryDeviceRepo.getById(id);
    this.validateEntryDeviceTransition(current.defunctReason, newDefunctReason);
    await this.entryDeviceRepo.update(id, deviceType, identifier, newDefunctReason);
  }

  async deleteEntryDevice(id: string): Promise<void> {
    await this.entryDeviceRepo.delete(id);
  }
}
