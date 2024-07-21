<script lang="ts">
  import DevlceList from '$lib/components/device/DevlceList.svelte';
  import { BackendClient } from '$lib/api/backend-client.svelte';
	import { getPropertyContext } from '$lib/contexts/property-context.svelte';
  
  const backendClient = new BackendClient();

  const propertyContext = getPropertyContext();
</script>
{propertyContext.selectedPropertyId}
<DevlceList backendClient={backendClient} propertyId={propertyContext.selectedPropertyId} />