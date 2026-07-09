import type { IPersonService } from '$lib/services/person/person-service';
import type { IPersonSetupModule } from '$lib/modules/interfaces';
import type { IPersonModel } from '$lib/models/person-models';

export class PersonSetupModule implements IPersonSetupModule {
	private readonly __personService: IPersonService;
	private readonly __userId: string;
	private readonly __backAction: () => void;

	public personAsync = $state<Promise<IPersonModel | null>>(
		Promise.resolve({ id: '', name: '', DOB: '', user: '', profileImage: '' })
	);

	private __error = $state<string>('');

	get isAdd() {
		return true;
	}

	public callBackAction = () => this.__backAction();

	public getSavePersonAction = () => async (person: IPersonModel) => {
		try {
			await this.__personService.createPerson(person.name, person.DOB, this.__userId);
			this.__backAction();
		} catch (ex) {
			this.__error = String(ex);
		}
	};

	public getDeletePersonAction = () => null;

	public get error() {
		return this.__error;
	}

	constructor(personService: IPersonService, userId: string, backAction: () => void) {
		this.__personService = personService;
		this.__userId = userId;
		this.__backAction = backAction;
	}
}
