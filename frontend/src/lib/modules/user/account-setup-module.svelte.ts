import type { IPersonService } from '$lib/services/person/person-service';
import type { ICobrandService } from '$lib/services/cobrand/cobrand-service';
import type { IAccountSetupModule } from '$lib/modules/interfaces';

export class AccountSetupModule implements IAccountSetupModule {
	public showChoiceAsync = $state<Promise<boolean>>(Promise.resolve(true));

	constructor(
		personService: IPersonService,
		cobrandService: ICobrandService,
		userId: string,
		redirectAction: (path: string) => void
	) {
		this.showChoiceAsync = (async () => {
			const person = await personService.getPersonByUserId(userId);
			if (person !== null) {
				redirectAction('/user/properties/list');
				return false;
			}
			const admin = await cobrandService.getAdminRecordForUserId(userId);
			if (admin !== null) {
				redirectAction('/user/cobrands/');
				return false;
			}
			return true;
		})();
	}
}
