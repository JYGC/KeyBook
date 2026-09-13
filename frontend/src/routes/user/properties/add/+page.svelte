<script lang="ts">
  import { goto } from '$app/navigation';
  import { getBackendClient } from '$lib/api/backend-client';
  import PropertyEditor from '$lib/components/property/PropertyEditor.svelte';
  import { PersonRepository } from '$lib/repositories/person/person-repository';
  import { PropertyRepository } from '$lib/repositories/property/property-repository';
  import { PropertyOwnerRepository } from '$lib/repositories/property/property-owner-repository';
  import { PersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
  import { TenantRepository } from '$lib/repositories/property/tenant-repository';
  import { HouseholdRepository } from '$lib/repositories/property/household-repository';
  import { PropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';
  import { PropertyService } from '$lib/services/property/property-service';
  import { PropertyAddModule } from '$lib/modules/property/property-add-module.svelte';
  import { Button, Tile } from 'carbon-components-svelte';

  const goBack = () => goto('/user/properties/list');

  const backendClient = getBackendClient();
  const personRepository = new PersonRepository(backendClient);
  const propertyRepository = new PropertyRepository(backendClient);
  const propertyOwnerRepository = new PropertyOwnerRepository(backendClient);
  const personPropertyOwnerRepository = new PersonPropertyOwnerRepository(backendClient);
  const tenantRepository = new TenantRepository(backendClient);
  const householdRepository = new HouseholdRepository(backendClient);
  const propertyAgentRepository = new PropertyAgentRepository(backendClient);
  const propertyService = new PropertyService(propertyRepository, propertyOwnerRepository, personPropertyOwnerRepository, tenantRepository, householdRepository, propertyAgentRepository);

  let propertyAddModule = $state<PropertyAddModule | null>(null);

  personRepository.getPersonByUserId(backendClient.authStore.record?.id ?? '').then((person) => {
    if (!person) {
      alert('No person record found for current user. Create a person first.');
      goBack();
      return;
    }
    propertyAddModule = new PropertyAddModule(propertyService, person.id, goBack);
  });
</script>

<Button onclick={goBack}>Back</Button>

{#if propertyAddModule}
  <PropertyEditor propertyEditorModule={propertyAddModule} />
{:else}
  <Tile>...loading</Tile>
{/if}
