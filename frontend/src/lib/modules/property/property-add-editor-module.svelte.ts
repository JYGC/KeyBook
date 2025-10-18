import type { IEditPropertyModel } from "$lib/models/property-models";
import PocketBase from "pocketbase";
import type { IPropertyEditorModule } from "../interfaces";

export class PropertyAddEditorModule implements IPropertyEditorModule {
  private readonly __backendClient: PocketBase;
  private readonly __backAction: () => void;
  
  constructor(
    backendClient: PocketBase,
    actionAfterSavePropertyAction: () => void,
  ) {
    this.__backendClient = backendClient;
    this.__backAction = actionAfterSavePropertyAction;
  }

  public propertyAsync = $state<Promise<IEditPropertyModel | null>>((async () => ({} as IEditPropertyModel))());
  
  get isAdd() { return true; }

  public getSavePropertyAction = () => (async (changedProperty: IEditPropertyModel) => {
    try {
      if (this.__backendClient.authStore.record === null) {
        throw new Error("Cannot find loggedInUser.");
      }
      changedProperty.owners = [ this.__backendClient.authStore.record.id ];
      await this.__backendClient.collection("properties").create<IEditPropertyModel>(changedProperty);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  });
  
  public getDeletePropertyAction = () => null;

  public callBackAction = () => this.__backAction();
}