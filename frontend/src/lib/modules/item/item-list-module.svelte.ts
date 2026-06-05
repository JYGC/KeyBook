import type { IItemService } from '$lib/services/item/item-service';
import type { IItemListModule } from '$lib/modules/interfaces';
import type { IItemModel } from '$lib/models/item-models';

export class ItemListModule implements IItemListModule {
  private readonly __itemService: IItemService;

  public itemListAsync = $derived.by<Promise<IItemModel[]>>(async () => {
    try {
      return await this.__itemService.getAllItems();
    } catch (ex) {
      alert(ex);
      return [];
    }
  });

  constructor(itemService: IItemService) {
    this.__itemService = itemService;
  }
}
