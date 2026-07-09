<script lang="ts">
  import { goto } from '$app/navigation';
  import { getBackendClient } from '$lib/api/backend-client';
  import PropertyList from '$lib/components/property/PropertyList.svelte';
  import { PropertyRepository } from '$lib/repositories/property/property-repository';
  import { PropertyOwnerRepository } from '$lib/repositories/property/property-owner-repository';
  import { PersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
  import { TenantRepository } from '$lib/repositories/property/tenant-repository';
  import { HouseholdRepository } from '$lib/repositories/property/household-repository';
  import { PropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';
  import { PropertyService } from '$lib/services/property/property-service';
  import { PropertyListModule } from '$lib/modules/property/property-list-module.svelte';
  import { Button } from 'carbon-components-svelte';

  const pb = getBackendClient();
  const propertyRepo = new PropertyRepository(pb);
  const propertyOwnerRepo = new PropertyOwnerRepository(pb);
  const ppoRepo = new PersonPropertyOwnerRepository(pb);
  const tenantRepo = new TenantRepository(pb);
  const householdRepo = new HouseholdRepository(pb);
  const propertyAgentRepo = new PropertyAgentRepository(pb);
  const propertyService = new PropertyService(propertyRepo, propertyOwnerRepo, ppoRepo, tenantRepo, householdRepo, propertyAgentRepo);
  const propertyListModule = new PropertyListModule(propertyService);

  const gotoAdd = () => goto('/user/properties/add');
</script>

<Button onclick={gotoAdd}>Add Property</Button>
<PropertyList {propertyListModule} />
