import PocketBase from "pocketbase";
import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IDeviceListModule } from "$lib/interfaces";
import type { IDeviceListItemModel } from "$lib/models/person-device-models";

export class DeviceListModule implements IDeviceListModule {
  private readonly __backendClient: PocketBase;
  private readonly __propertyContext: PropertyContext;

  public deviceListAsync = $derived.by<Promise<IDeviceListItemModel[]>>(async () => {
    try {
      return await this.__backendClient.collection("devicelistview").getFullList<IDeviceListItemModel>({
        filter: `propertyid = "${this.__propertyContext.selectedPropertyId}"`,
        fields: "id,deviceid,devicename,deviceidentifier,devicetype,personname",
      });
    } catch (ex) {
      alert(ex);
      return [];
    }
  });
  
  constructor(backendClient: PocketBase, propertyContext: PropertyContext) {
    this.__backendClient = backendClient;
    this.__propertyContext = propertyContext;
  }
}