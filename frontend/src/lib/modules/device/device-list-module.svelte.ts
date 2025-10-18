import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IDeviceListModule } from "$lib/interfaces";
import type { IDeviceListItemModel } from "$lib/models/person-device-models";
import type { IDeviceListViewService } from "$lib/services/interfaces";

export class DeviceListModule implements IDeviceListModule {
  private readonly __deviceListViewService: IDeviceListViewService;
  private readonly __propertyContext: PropertyContext;

  public deviceListAsync = $derived.by<Promise<IDeviceListItemModel[]>>(async () => {
    try {
      return await this.__deviceListViewService.getForDeviceListAsync(this.__propertyContext.selectedPropertyId);
    } catch (ex) {
      alert(ex);
      return [];
    }
  });
  
  constructor(deviceListViewService: IDeviceListViewService, propertyContext: PropertyContext) {
    this.__deviceListViewService = deviceListViewService;
    this.__propertyContext = propertyContext;
  }
}