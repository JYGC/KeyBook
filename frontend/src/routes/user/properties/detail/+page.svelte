<script lang="ts">
  import { goto } from '$app/navigation';
  import { getBackendClient } from '$lib/api/backend-client';
  import PropertyEditor from '$lib/components/property/PropertyEditor.svelte';
  import PersonOwnerList from '$lib/components/property/PersonOwnerList.svelte';
  import TenantList from '$lib/components/property/TenantList.svelte';
  import HouseholdMemberList from '$lib/components/property/HouseholdMemberList.svelte';
  import ConfirmButtonAndDialog from '$lib/components/shared/ConfirmButtonAndDialog.svelte';
  import { PropertyRepository } from '$lib/repositories/property/property-repository';
  import { PropertyOwnerRepository } from '$lib/repositories/property/property-owner-repository';
  import { PersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
  import { TenantRepository } from '$lib/repositories/property/tenant-repository';
  import { HouseholdRepository } from '$lib/repositories/property/household-repository';
  import { PropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';
  import { PropertyService } from '$lib/services/property/property-service';
  import { PropertyDetailModule } from '$lib/modules/property/property-detail-module.svelte';
  import { Button } from 'carbon-components-svelte';

  let { data } = $props();
  const propertyId: string = data.propertyId;

  const goBack = () => goto('/user/properties/list');

  const pb = getBackendClient();
  const propertyRepo = new PropertyRepository(pb);
  const propertyOwnerRepo = new PropertyOwnerRepository(pb);
  const ppoRepo = new PersonPropertyOwnerRepository(pb);
  const tenantRepo = new TenantRepository(pb);
  const householdRepo = new HouseholdRepository(pb);
  const propertyAgentRepo = new PropertyAgentRepository(pb);
  const propertyService = new PropertyService(propertyRepo, propertyOwnerRepo, ppoRepo, tenantRepo, householdRepo, propertyAgentRepo);
  const propertyDetailModule = new PropertyDetailModule(propertyService, propertyId, goBack);
</script>

<Button onclick={goBack}>Back</Button>

<PropertyEditor propertyEditorModule={propertyDetailModule}>
  {#snippet deleteButton(deleteActionButtonClick: () => void)}
    <ConfirmButtonAndDialog
      submitAction={deleteActionButtonClick}
      buttonText="Delete Property"
      bodyMessage="Are you sure you want to delete this property?"
    />
  {/snippet}
</PropertyEditor>

<br />

<PersonOwnerList {propertyDetailModule} />

<br />

<TenantList {propertyId} {propertyDetailModule} />

<br />

<HouseholdMemberList {propertyId} {propertyDetailModule} />
