import type { IPersonService } from '$lib/services/person/person-service';
import type { INewPersonListModule } from '$lib/modules/interfaces';
import type { IPersonModel } from '$lib/models/person-models';

export class PersonListModule implements INewPersonListModule {
  private readonly __personService: IPersonService;

  public personListAsync = $derived.by<Promise<IPersonModel[]>>(async () => {
    try {
      return await this.__personService.getAllPersons();
    } catch (ex) {
      alert(ex);
      return [];
    }
  });

  constructor(personService: IPersonService) {
    this.__personService = personService;
  }
}
