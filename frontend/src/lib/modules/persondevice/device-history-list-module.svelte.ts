import type { DeviceContext } from "$lib/contexts/device-context.svelte";
import type { IBackendClient, IDeviceHistoryListModule } from "$lib/interfaces";
import type { IDeviceHistoryListItem as IDeviceHistoryListItemModel } from "$lib/models/person-device-models";

export class DeviceHistoryListModule implements IDeviceHistoryListModule {
  private readonly __backendClient: IBackendClient;
  private readonly __deviceContext: DeviceContext;

  public deviceHistoryListAsync = $derived.by<Promise<IDeviceHistoryListItemModel[]>>(async () => {
    try {
      const ite =  await this.__backendClient.pb.collection("devicehistorylistview").getFullList<IDeviceHistoryListItemModel>({
        filter: `deviceid = "${this.__deviceContext.selectedDeviceId}"`,
        fields: "id,deviceid,description,created",
      });
      console.log(ite);
      return ite;
    } catch (ex) {
      alert(ex);
      return []
    }
  });

  constructor(
    backendClient: IBackendClient,
    deviceContext: DeviceContext
  ) {
    this.__backendClient = backendClient;
    this.__deviceContext = deviceContext;
  }
}