import type { IPersonService } from '$lib/services/person/person-service';
import type { INewPersonEditorModule } from '$lib/modules/interfaces';
import type { IPersonModel } from '$lib/models/person-models';

export class PersonAddModule implements INewPersonEditorModule {
  private readonly __personService: IPersonService;
  private readonly __backAction: () => void;

  public personAsync = $state<Promise<IPersonModel | null>>(
    Promise.resolve({ id: '', name: '', DOB: '', user: '', profileImage: '' }),
  );

  get isAdd() {
    return true;
  }

  public callBackAction = () => this.__backAction();

  public getSavePersonAction = () => async (person: IPersonModel) => {
    try {
      await this.__personService.createPerson(person.name, person.DOB);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getDeletePersonAction = () => null;

  constructor(personService: IPersonService, backAction: () => void) {
    this.__personService = personService;
    this.__backAction = backAction;
  }
}
