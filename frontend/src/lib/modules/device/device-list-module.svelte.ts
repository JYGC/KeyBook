import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IBackendClient, IDeviceListModule } from "$lib/interfaces";
import type { IDeviceListItemModel } from "$lib/models/person-device-models";

export class DeviceListModule implements IDeviceListModule {
  private readonly __backendClient: IBackendClient;
  private readonly __propertyContext: PropertyContext;

  public deviceListAsync = $derived.by<Promise<IDeviceListItemModel[]>>(async () => {
    try {
      const items = await this.__backendClient.pb.collection("devicelistview").getFullList<IDeviceListItemModel>({
        filter: `propertyid = "${this.__propertyContext.selectedPropertyId}"`,
        fields: "id,deviceid,devicename,deviceidentifier,devicetype,personname",
      });
      return items;
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