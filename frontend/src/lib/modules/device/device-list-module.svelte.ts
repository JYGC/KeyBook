import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IBackendClient, IDeviceListModule } from "$lib/interfaces";
import type { IDeviceListItemModel } from "$lib/models/device-models";

export class DeviceListModule implements IDeviceListModule {
  private readonly __backendClient: IBackendClient;
  private readonly __propertyContext: PropertyContext;

  public deviceListAsync = $derived.by<Promise<IDeviceListItemModel[]>>(async () => {
    try {
      const response = await this.__backendClient.pb.collection("devices").getList<IDeviceListItemModel>(1, 50, {
        filter: `property.id = "${this.__propertyContext.selectedPropertyId}"`,
        fields: "id,type,name,identifier",
      });
      return response.items;
    } catch (ex) {
      alert(ex);
      return [];
    }
  });
  
  constructor(backendClient: IBackendClient, propertyContext: PropertyContext) {
    this.__backendClient = backendClient;
    this.__propertyContext = propertyContext;
  }
}