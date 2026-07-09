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

  const pb = getBackendClient();
  const personRepo = new PersonRepository(pb);
  const propertyRepo = new PropertyRepository(pb);
  const propertyOwnerRepo = new PropertyOwnerRepository(pb);
  const ppoRepo = new PersonPropertyOwnerRepository(pb);
  const tenantRepo = new TenantRepository(pb);
  const householdRepo = new HouseholdRepository(pb);
  const propertyAgentRepo = new PropertyAgentRepository(pb);
  const propertyService = new PropertyService(propertyRepo, propertyOwnerRepo, ppoRepo, tenantRepo, householdRepo, propertyAgentRepo);

  let propertyAddModule = $state<PropertyAddModule | null>(null);

  personRepo.getByUserId(pb.authStore.record?.id ?? '').then((person) => {
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
