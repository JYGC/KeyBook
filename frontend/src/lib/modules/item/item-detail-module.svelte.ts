import type { IItemService } from '$lib/services/item/item-service';
import type { IItemEditorModule } from '$lib/modules/interfaces';
import type { IItemModel, IEntryDeviceModel } from '$lib/models/item-models';

export class ItemDetailModule implements IItemEditorModule {
  private readonly __itemService: IItemService;
  private readonly __itemId: string;
  private readonly __backAction: () => void;

  public itemAsync: Promise<IItemModel | null>;
  public entryDeviceAsync: Promise<IEntryDeviceModel | null>;

  get isAdd() {
    return false;
  }

  public callBackAction = () => this.__backAction();

  public getSaveItemAction = () => async (item: IItemModel) => {
    try {
      await this.__itemService.updateItem(item.id, item.name, item.description);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getDeleteItemAction = () => async (id: string) => {
    try {
      await this.__itemService.deleteItem(id);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getSaveEntryDeviceAction = () => async (ed: IEntryDeviceModel) => {
    try {
      await this.__itemService.updateEntryDevice(ed.id, ed.deviceType, ed.identifier, ed.defunctReason);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getDeleteEntryDeviceAction = () => async (id: string) => {
    try {
      await this.__itemService.deleteEntryDevice(id);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getCreateEntryDeviceAction =
    () => async (itemId: string, deviceType: string, identifier: string, defunctReason: string) => {
      try {
        await this.__itemService.createEntryDevice(itemId, deviceType, identifier, defunctReason);
        this.__backAction();
      } catch (ex) {
        alert(ex);
      }
    };

  constructor(itemService: IItemService, itemId: string, backAction: () => void) {
    this.__itemService = itemService;
    this.__itemId = itemId;
    this.__backAction = backAction;
    const combined = this.__itemService.getItemWithEntryDevice(this.__itemId)
      .catch((ex) => { alert(ex); return null; });
    this.itemAsync = combined.then(result => result?.item ?? null);
    this.entryDeviceAsync = combined.then(result => result?.entryDevice ?? null);
  }
}
