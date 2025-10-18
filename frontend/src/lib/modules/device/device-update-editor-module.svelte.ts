import PocketBase from "pocketbase";
import type { DeviceContext } from "$lib/contexts/device-context.svelte";
import type { IEditDeviceModel } from "$lib/models/device-models";
import type { IDeviceEditorModule } from "$lib/modules/interfaces";

export class DeviceUpdateEditorModule implements IDeviceEditorModule {
  private readonly __backendClient: PocketBase;
  private readonly __deviceContext: DeviceContext;
  private readonly __backAction: () => void;

  public deviceAsync = $derived.by<Promise<IEditDeviceModel | null>>(async () => {
    try {
      return await this.__backendClient.collection("devices").getOne<IEditDeviceModel>(
        this.__deviceContext.selectedDeviceId,
        { fields: "id,type,name,identifier,defunctreason,property" },
      );
    } catch (ex) {
      alert(ex);
      return null;
    }
  });

  get isAdd() { return false; }

  public deviceStatusTextAsync = $derived.by<Promise<string>>(async () => {
    const device = await this.deviceAsync;
    return (device === null || device.defunctreason === "None") ? "Usable" : device.defunctreason;
  });

  public getSaveDeviceAction = () => (async (changedDevice: IEditDeviceModel) => {
    try {
      await this.__backendClient.collection("devices").update(changedDevice.id, {
        type: changedDevice.type,
        name: changedDevice.name,
        identifier: changedDevice.identifier,
        defunctreason: changedDevice.defunctreason,
      });
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  });

  public getDeleteDeviceAction = () => (async (device: IEditDeviceModel) => {
    try {
      await this.__backendClient.collection("devices").delete(device.id);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  });

  public callBackAction = () => this.__backAction();
  
  constructor(
    backendClient: PocketBase,
    deviceContext: DeviceContext,
    actionAfterSaveDeviceAction: () => void,
  ) {
    this.__backendClient = backendClient;
    this.__deviceContext = deviceContext;
    this.__backAction = actionAfterSaveDeviceAction;
  }
}