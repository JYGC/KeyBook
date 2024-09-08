import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IEditDeviceDto } from "$lib/dtos/device-dtos";
import type { IBackendClient, IDeviceEditorModule } from "$lib/interfaces";

export class DeviceAddEditorModule implements IDeviceEditorModule {
  private readonly __backendClient: IBackendClient;
  private readonly __propertyContext: PropertyContext;
  private readonly __backAction: () => void;

  public deviceAsync = $state<Promise<IEditDeviceDto | null>>((async () => ({} as IEditDeviceDto))());
  
  get isAdd() { return true; }
  
  public deviceStatusTextAsync = (async () => "")();

  public getSaveDeviceAction = () => (async (changedDevice: IEditDeviceDto) => {
    try {
      changedDevice.property = this.__propertyContext.selectedPropertyId;
      await this.__backendClient.pb.collection("devices").create<IEditDeviceDto>(changedDevice);
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