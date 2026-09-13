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
  import { PersonAddModule } from '$lib/modules/person/person-add-module.svelte';
  import { Button } from 'carbon-components-svelte';

  const goBack = () => goto('/user/persons');

  const backendClient = getBackendClient();
  const personRepository = new PersonRepository(backendClient);
  const personPropertyOwnerRepository = new PersonPropertyOwnerRepository(backendClient);
  const tenantRepository = new TenantRepository(backendClient);
  const householdRepository = new HouseholdRepository(backendClient);
  const agentRepository = new AgentRepository(backendClient);
  const personService = new PersonService(personRepository, personPropertyOwnerRepository, tenantRepository, householdRepository, agentRepository);
  const personEditorModule = new PersonAddModule(personService, goBack);
</script>

<Button onclick={goBack}>Back</Button>
<PersonEditor {personEditorModule} />
