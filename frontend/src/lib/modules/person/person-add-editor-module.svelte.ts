import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IBackendClient, IPersonEditorModule } from "$lib/interfaces";
import type { IEditPersonModel } from "$lib/models/person-models";

export class PersonAddEditorModule implements IPersonEditorModule {
  private readonly __backendClient: IBackendClient;
  private readonly __propertyContext: PropertyContext;
  private readonly __backAction: () => void;

  public personAsync = $state<Promise<IEditPersonModel | null>>((async () => ({} as IEditPersonModel))());

  public get isAdd(): boolean { return true; }
  
  public getSavePersonAction = () => (async (changedPerson: IEditPersonModel) => {
    try {
      changedPerson.property = this.__propertyContext.selectedPropertyId;
      await this.__backendClient.pb.collection("persons").create<IEditPersonModel>(changedPerson);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  });
  public getDeletePersonAction = () => null;
  public callBackAction = () => this.__backAction();

  constructor(
    backendClient: IBackendClient,
    personContext: PropertyContext,
    actionAfterSavePersonAction: () => void,
  ) {
    this.__backendClient = backendClient;
    this.__propertyContext = personContext;
    this.__backAction = actionAfterSavePersonAction;
  }
}