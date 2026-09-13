import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { getBackendClient } from '$lib/api/backend-client';
import { PersonRepository } from '$lib/repositories/person/person-repository';
import { PersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
import { TenantRepository } from '$lib/repositories/property/tenant-repository';
import { HouseholdRepository } from '$lib/repositories/property/household-repository';
import { AgentRepository } from '$lib/repositories/agent/agent-repository';
import { PersonService } from '$lib/services/person/person-service';
import { CobrandRepository } from '$lib/repositories/cobrand/cobrand-repository';
import { CobrandAdminRepository } from '$lib/repositories/cobrand/cobrand-admin-repository';
import { CobrandPropertyManagerRepository } from '$lib/repositories/cobrand/cobrand-property-manager-repository';
import { CobrandService } from '$lib/services/cobrand/cobrand-service';

export const load: PageLoad = async () => {
	const backendClient = getBackendClient();
	const userId = backendClient.authStore.record?.id ?? '';

	const personService = new PersonService(
		new PersonRepository(backendClient),
		new PersonPropertyOwnerRepository(backendClient),
		new TenantRepository(backendClient),
		new HouseholdRepository(backendClient),
		new AgentRepository(backendClient)
	);
	const person = await personService.getPersonByUserId(userId);
	if (person !== null) {
		return redirect(303, '/user/properties/list');
	}

	const cobrandService = new CobrandService(
		new CobrandRepository(backendClient),
		new CobrandAdminRepository(backendClient),
		new CobrandPropertyManagerRepository(backendClient)
	);
	const admin = await cobrandService.getAdminRecordForUserId(userId);
	if (admin !== null) {
		return redirect(303, '/user/cobrands/');
	}

	return redirect(303, '/user/setup');
};
