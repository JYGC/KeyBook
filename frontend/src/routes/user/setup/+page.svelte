<script lang="ts">
	import { goto } from '$app/navigation';
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
	import { AccountSetupModule } from '$lib/modules/user/account-setup-module.svelte';
	import { Button, Tile } from 'carbon-components-svelte';

	const backendClient = getBackendClient();
	const currentUserId = backendClient.authStore.record?.id ?? '';

	const personService = new PersonService(
		new PersonRepository(backendClient),
		new PersonPropertyOwnerRepository(backendClient),
		new TenantRepository(backendClient),
		new HouseholdRepository(backendClient),
		new AgentRepository(backendClient)
	);
	const cobrandService = new CobrandService(
		new CobrandRepository(backendClient),
		new CobrandAdminRepository(backendClient),
		new CobrandPropertyManagerRepository(backendClient)
	);
	const accountSetupModule = new AccountSetupModule(
		personService,
		cobrandService,
		currentUserId,
		(path) => goto(path)
	);
</script>

{#await accountSetupModule.showChoiceAsync}
	<Tile>...checking your account</Tile>
{:then showChoice}
	{#if showChoice}
		<h1>Set up your account</h1>
		<br />
		<p>Choose how you'll be using KeyBook.</p>
		<br />
		<Button onclick={() => goto('/user/persons/setup/')}>Set up my person profile</Button>
		<Button onclick={() => goto('/user/cobrands/add/')}>Set up my company</Button>
	{/if}
{/await}
