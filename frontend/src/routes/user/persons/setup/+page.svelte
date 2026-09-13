<script lang="ts">
	import { goto } from '$app/navigation';
	import { getBackendClient } from '$lib/api/backend-client';
	import PersonEditor from '$lib/components/person/PersonEditor.svelte';
	import { PersonRepository } from '$lib/repositories/person/person-repository';
	import { PersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
	import { TenantRepository } from '$lib/repositories/property/tenant-repository';
	import { HouseholdRepository } from '$lib/repositories/property/household-repository';
	import { AgentRepository } from '$lib/repositories/agent/agent-repository';
	import { PersonService } from '$lib/services/person/person-service';
	import { PersonSetupModule } from '$lib/modules/person/person-setup-module.svelte';

	const goToProperties = () => goto('/user/properties/list');

	const backendClient = getBackendClient();
	const currentUserId = backendClient.authStore.record?.id ?? '';
	const personRepository = new PersonRepository(backendClient);
	const personPropertyOwnerRepository = new PersonPropertyOwnerRepository(backendClient);
	const tenantRepository = new TenantRepository(backendClient);
	const householdRepository = new HouseholdRepository(backendClient);
	const agentRepository = new AgentRepository(backendClient);
	const personService = new PersonService(
		personRepository,
		personPropertyOwnerRepository,
		tenantRepository,
		householdRepository,
		agentRepository
	);
	const personSetupModule = new PersonSetupModule(personService, currentUserId, goToProperties);
</script>

<h1>Set up your person profile</h1>
<br />
<PersonEditor personEditorModule={personSetupModule} />
{#if personSetupModule.error}
	<p class="error">{personSetupModule.error}</p>
{/if}
