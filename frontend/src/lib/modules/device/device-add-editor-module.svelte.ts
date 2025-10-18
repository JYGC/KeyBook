import PocketBase from "pocketbase";
import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IEditDeviceModel } from "$lib/models/device-models";
import type { IDeviceEditorModule } from "$lib/modules/interfaces";

export class DeviceAddEditorModule implements IDeviceEditorModule {
  private readonly __backendClient: PocketBase;
  private readonly __propertyContext: PropertyContext;
  private readonly __backAction: () => void;

  public deviceAsync = $state<Promise<IEditDeviceModel | null>>((async () => ({} as IEditDeviceModel))());
  
  get isAdd() { return true; }
  
  public deviceStatusTextAsync = (async () => "")();

  public getSaveDeviceAction = () => (async (changedDevice: IEditDeviceModel) => {
    try {
      changedDevice.property = this.__propertyContext.selectedPropertyId;
      await this.__backendClient.collection("devices").create<IEditDeviceModel>(changedDevice);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  });
  
  public getDeleteDeviceAction = () => null;

  public callBackAction = () => this.__backAction();
  
  constructor(
    backendClient: PocketBase,
    propertyContext: PropertyContext,
    actionAfterSaveDeviceAction: () => void,
  ) {
    this.__backendClient = backendClient;
    this.__propertyContext = propertyContext;
    this.__backAction = actionAfterSaveDeviceAction;
  }
}