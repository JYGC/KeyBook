<script lang="ts">
  import { goto } from '$app/navigation';
  import { getBackendClient } from '$lib/api/backend-client';
  import PersonList from '$lib/components/person/PersonList.svelte';
  import { PersonRepository } from '$lib/repositories/person/person-repository';
  import { PersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
  import { TenantRepository } from '$lib/repositories/property/tenant-repository';
  import { HouseholdRepository } from '$lib/repositories/property/household-repository';
  import { AgentRepository } from '$lib/repositories/agent/agent-repository';
  import { PersonService } from '$lib/services/person/person-service';
  import { PersonListModule } from '$lib/modules/person/person-list-module.svelte';
  import { Button } from 'carbon-components-svelte';

  const pb = getBackendClient();
  const personRepo = new PersonRepository(pb);
  const ppoRepo = new PersonPropertyOwnerRepository(pb);
  const tenantRepo = new TenantRepository(pb);
  const householdRepo = new HouseholdRepository(pb);
  const agentRepo = new AgentRepository(pb);
  const personService = new PersonService(personRepo, ppoRepo, tenantRepo, householdRepo, agentRepo);
  const personListModule = new PersonListModule(personService);

  const gotoAdd = () => goto('/user/persons/add');
</script>

<Button onclick={gotoAdd}>Add Person</Button>
<PersonList {personListModule} />
