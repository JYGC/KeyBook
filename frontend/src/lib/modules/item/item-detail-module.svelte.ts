import type { IItemService } from '$lib/services/item/item-service';
import type { ItemContext } from '$lib/contexts/item-context.svelte';
import type { IItemEditorModule } from '$lib/modules/interfaces';
import type { IItemModel, IEntryDeviceModel } from '$lib/models/item-models';

export class ItemDetailModule implements IItemEditorModule {
  private readonly __itemService: IItemService;
  private readonly __context: ItemContext;
  private readonly __backAction: () => void;

  public itemAsync = $derived.by<Promise<IItemModel | null>>(async () => {
    try {
      const { item } = await this.__itemService.getItemWithEntryDevice(this.__context.selectedItemId);
      return item;
    } catch (ex) {
      alert(ex);
      return null;
    }
  });

  public entryDeviceAsync = $derived.by<Promise<IEntryDeviceModel | null>>(async () => {
    try {
      const { entryDevice } = await this.__itemService.getItemWithEntryDevice(
        this.__context.selectedItemId,
      );
      return entryDevice;
    } catch (ex) {
      alert(ex);
      return null;
    }
  });

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

  constructor(itemService: IItemService, context: ItemContext, backAction: () => void) {
    this.__itemService = itemService;
    this.__context = context;
    this.__backAction = backAction;
  }
}
