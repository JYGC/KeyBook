import type { DeviceContext } from "$lib/contexts/device-context.svelte";
import type { IEditDeviceDto } from "$lib/dtos/device-dtos";
import type { IBackendClient, IDeviceEditorModule } from "$lib/interfaces";

export class DeviceUpdateEditorModule implements IDeviceEditorModule {
  private readonly __backendClient: IBackendClient;
  private readonly __deviceContext: DeviceContext;
  private readonly __backAction: () => void;

  public deviceAsync = $derived.by<Promise<IEditDeviceDto | null>>(async () => {
    try {
      return await this.__backendClient.pb.collection("devices").getOne<IEditDeviceDto>(
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

  public getSaveDeviceAction = () => (async (changedDevice: IEditDeviceDto) => {
    try {
      await this.__backendClient.pb.collection("devices").update(changedDevice.id, {
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

  public getDeleteDeviceAction = () => (async (device: IEditDeviceDto) => {
    try {
      await this.__backendClient.pb.collection("devices").delete(device.id);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  });

  public callBackAction = () => this.__backAction();
  
  constructor(
    backendClient: IBackendClient,
    deviceContext: DeviceContext,
    actionAfterSaveDeviceAction: () => void,
  ) {
    this.__backendClient = backendClient;
    this.__deviceContext = deviceContext;
    this.__backAction = actionAfterSaveDeviceAction;
  }
}