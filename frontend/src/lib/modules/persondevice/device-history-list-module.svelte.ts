import PocketBase from "pocketbase";
import type { DeviceContext } from "$lib/contexts/device-context.svelte";
import type { IDeviceHistoryListModule, IDeviceHistoryListUpdaterModule } from "$lib/modules/nterfaces";
import type { IDeviceHistoryListItem as IDeviceHistoryListItemModel } from "$lib/models/person-device-models";

export class DeviceHistoryListModule implements IDeviceHistoryListModule {
  private readonly __backendClient: PocketBase;
  private readonly __deviceContext: DeviceContext;
  private readonly __deviceHistoryListUpdaterModule: IDeviceHistoryListUpdaterModule | null;

  private previousUpdateTriggerState: unknown;

  public deviceHistoryListAsync = $derived.by<Promise<IDeviceHistoryListItemModel[]>>(async () => {
    try {
      if (this.__deviceHistoryListUpdaterModule !== null && this.previousUpdateTriggerState !== this.__deviceHistoryListUpdaterModule.updateTriggerState) {
        this.previousUpdateTriggerState = this.__deviceHistoryListUpdaterModule.updateTriggerState;
      }
      return (await this.__backendClient.collection("deviceHistoryListView").getFullList<IDeviceHistoryListItemModel>({
        filter: `deviceId = "${this.__deviceContext.selectedDeviceId}"`,
        fields: "id,deviceId,description,created",
      })).sort((a, b) => new Date(b.created) - new Date(a.created));
    } catch (ex) {
      alert(ex);
      return []
    }
  });

  constructor(
    backendClient: PocketBase,
    deviceContext: DeviceContext,
    deviceHistoryListUpdaterModule: IDeviceHistoryListUpdaterModule | null = null
  ) {
    this.__backendClient = backendClient;
    this.__deviceContext = deviceContext;
    this.__deviceHistoryListUpdaterModule = deviceHistoryListUpdaterModule;
  }
}
