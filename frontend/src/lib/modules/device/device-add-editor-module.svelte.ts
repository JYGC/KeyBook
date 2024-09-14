import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IEditDeviceModel } from "$lib/dtos/device-models";
import type { IBackendClient, IDeviceEditorModule } from "$lib/interfaces";

export class DeviceAddEditorModule implements IDeviceEditorModule {
  private readonly __backendClient: IBackendClient;
  private readonly __propertyContext: PropertyContext;
  private readonly __backAction: () => void;

  public deviceAsync = $state<Promise<IEditDeviceModel | null>>((async () => ({} as IEditDeviceModel))());
  
  get isAdd() { return true; }
  
  public deviceStatusTextAsync = (async () => "")();

  public getSaveDeviceAction = () => (async (changedDevice: IEditDeviceModel) => {
    try {
      changedDevice.property = this.__propertyContext.selectedPropertyId;
      await this.__backendClient.pb.collection("devices").create<IEditDeviceModel>(changedDevice);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  });
  
  public getDeleteDeviceAction = () => null;

  public callBackAction = () => this.__backAction();
  
  constructor(
    backendClient: IBackendClient,
    propertyContext: PropertyContext,
    actionAfterSaveDeviceAction: () => void,
  ) {
    this.__backendClient = backendClient;
    this.__propertyContext = propertyContext;
    this.__backAction = actionAfterSaveDeviceAction;
  }
}