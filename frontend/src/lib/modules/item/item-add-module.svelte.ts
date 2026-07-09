import type { IItemService } from '$lib/services/item/item-service';
import type { IItemEditorModule } from '$lib/modules/interfaces';
import type { IItemModel, IEntryDeviceModel } from '$lib/models/item-models';

export class ItemAddModule implements IItemEditorModule {
  private readonly __itemService: IItemService;
  private readonly __backAction: () => void;

  public itemAsync = $state<Promise<IItemModel | null>>(
    Promise.resolve({ id: '', name: '', description: '' }),
  );

  public entryDeviceAsync = $state<Promise<IEntryDeviceModel | null>>(Promise.resolve(null));

  get isAdd() {
    return true;
  }

  public callBackAction = () => this.__backAction();

  public getSaveItemAction = () => async (item: IItemModel) => {
    try {
      await this.__itemService.createItem(item.name, item.description);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getDeleteItemAction = () => null;

  constructor(itemService: IItemService, backAction: () => void) {
    this.__itemService = itemService;
    this.__backAction = backAction;
  }
}
