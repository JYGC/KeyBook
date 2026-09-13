<script lang="ts">
  import { goto } from '$app/navigation';
  import { getBackendClient } from '$lib/api/backend-client';
  import PersonEditor from '$lib/components/person/PersonEditor.svelte';
  import ConfirmButtonAndDialog from '$lib/components/shared/ConfirmButtonAndDialog.svelte';
  import { PersonRepository } from '$lib/repositories/person/person-repository';
  import { PersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
  import { TenantRepository } from '$lib/repositories/property/tenant-repository';
  import { HouseholdRepository } from '$lib/repositories/property/household-repository';
  import { AgentRepository } from '$lib/repositories/agent/agent-repository';
  import { PersonService } from '$lib/services/person/person-service';
  import { PersonDetailModule } from '$lib/modules/person/person-detail-module.svelte';
  import { Button, Tile } from 'carbon-components-svelte';

  let { data } = $props();
  const personId: string = data.personId;

  const goBack = () => goto('/user/persons');

  const backendClient = getBackendClient();
  const personRepository = new PersonRepository(backendClient);
  const personPropertyOwnerRepository = new PersonPropertyOwnerRepository(backendClient);
  const tenantRepository = new TenantRepository(backendClient);
  const householdRepository = new HouseholdRepository(backendClient);
  const agentRepository = new AgentRepository(backendClient);
  const personService = new PersonService(personRepository, personPropertyOwnerRepository, tenantRepository, householdRepository, agentRepository);
  const personDetailModule = new PersonDetailModule(personService, personId, goBack);
</script>

<Button onclick={goBack}>Back</Button>

<PersonEditor personEditorModule={personDetailModule}>
  {#snippet deleteButton(deleteActionButtonClick: () => void)}
    <ConfirmButtonAndDialog
      submitAction={deleteActionButtonClick}
      buttonText="Delete Person"
      bodyMessage="Are you sure you want to delete this person?"
    />
  {/snippet}
</PersonEditor>

<br />

<h3>Roles</h3>
{#await personDetailModule.rolesAsync}
  <Tile>...getting roles</Tile>
{:then roles}
  {#if roles.length === 0}
    <p>No roles assigned.</p>
  {:else}
    <ul>
      {#each roles as role}
        <li>{role}</li>
      {/each}
    </ul>
  {/if}
{:catch error}
  {error}
{/await}
