import type { IPersonService } from '$lib/services/person/person-service';
import type { IPersonDetailModule } from '$lib/modules/interfaces';
import type { IPersonModel } from '$lib/models/person-models';

export class PersonDetailModule implements IPersonDetailModule {
  private readonly __personService: IPersonService;
  private readonly __personId: string;
  private readonly __backAction: () => void;

  public personAsync: Promise<IPersonModel | null>;
  public rolesAsync = $state<Promise<string[]>>(Promise.resolve([]));

  get isAdd() {
    return false;
  }

  public callBackAction = () => this.__backAction();

  private refreshRoles() {
    this.rolesAsync = this.__personService.getRolesForPerson(this.__personId).catch((ex) => {
      alert(ex);
      return [];
    });
  }

  public getSavePersonAction = () => async (person: IPersonModel) => {
    try {
      await this.__personService.updatePerson(person.id, person.name, person.DOB);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getDeletePersonAction = () => async (id: string) => {
    try {
      await this.__personService.deletePerson(id);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  constructor(personService: IPersonService, personId: string, backAction: () => void) {
    this.__personService = personService;
    this.__personId = personId;
    this.__backAction = backAction;

    this.personAsync = this.__personService.getPersonById(personId).catch((ex) => {
      alert(ex);
      return null;
    });

    this.refreshRoles();
  }
}
