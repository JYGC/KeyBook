import PocketBase from "pocketbase";
import type { PersonContext } from "$lib/contexts/person-context.svelte";
import type { IPersonEditorModule } from "$lib/interfaces";
import type { IEditPersonModel } from "$lib/models/person-models";

export class PersonUpdateEditorModule implements IPersonEditorModule {
  private readonly __backendClient: PocketBase;
  private readonly __personContext: PersonContext;
  private readonly __backAction: () => void;

  public personAsync = $derived.by<Promise<IEditPersonModel | null>>(async () => {
    try {
      return await this.__backendClient.collection("persons").getOne<IEditPersonModel>(
        this.__personContext.selectedPersonId,
        { fields: "id,name,type,property" },
      );
    } catch (ex) {
      alert(ex);
      return null;
    }
  });

  public get isAdd(): boolean { return false; }

  public getSavePersonAction = () => (async (changedPerson: IEditPersonModel) => {
    try {
      await this.__backendClient.collection("persons").update(changedPerson.id, {
        type: changedPerson.type,
        name: changedPerson.name,
      });
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  });
  
  public getDeletePersonAction = () => (async (person: IEditPersonModel) => {
    try {
      await this.__backendClient.collection("persons").delete(person.id);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  });
  
  public callBackAction = () => this.__backAction();

  constructor(
    backendClient: PocketBase,
    personContext: PersonContext,
    actionAfterSavePersonAction: () => void,
  ) {
    this.__backendClient = backendClient;
    this.__personContext = personContext;
    this.__backAction = actionAfterSavePersonAction;
  }
}