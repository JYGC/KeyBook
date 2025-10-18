import type { IEditPropertyModel } from "$lib/models/property-models";
import PocketBase from "pocketbase";
import type { IPropertyEditorModule } from "../interfaces";
import type { PropertyContext } from "$lib/contexts/property-context.svelte";

export class PropertyUpdateEditorModule implements IPropertyEditorModule {
  private readonly __backendClient: PocketBase;
  private readonly __propertyContext: PropertyContext;
  private readonly __backAction: () => void;
  
  constructor(
    backendClient: PocketBase,
    propertyContext: PropertyContext,
    actionAfterSavePropertyAction: () => void,
  ) {
    this.__backendClient = backendClient;
    this.__propertyContext = propertyContext;
    this.__backAction = actionAfterSavePropertyAction;
  }

  public propertyAsync = $derived.by<Promise<IEditPropertyModel | null>>(async () => {
    try {
      return await this.__backendClient.collection("properties").getOne<IEditPropertyModel>(
        this.__propertyContext.selectedPropertyId,
        { fields: "id,address,owners,managers" },
      );
    } catch (ex) {
      alert(ex);
      return null;
    }
  });

  get isAdd() { return false; }

  public savePropertyAction = () => async (changedProperty: IEditPropertyModel) => {
    try {
      await this.__backendClient.collection("properties").update(changedProperty.id, {
        address: changedProperty.address,
        owners: changedProperty.owners,
        managers: changedProperty.managers,
      });
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };
  
  public getDeletePropertyAction = () => async (property: IEditPropertyModel) => {
    try {
      await this.__backendClient.collection("properties").delete(property.id);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };;

  public callBackAction = () => this.__backAction();
}