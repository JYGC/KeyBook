<script lang="ts">
  import DevlceList from '$lib/components/device/DevlceList.svelte';
  import { BackendClient } from '$lib/api/backend-client.svelte';
	import { PropertyContext } from '$lib/stores/property-context.svelte';
	import { getContext } from 'svelte';
  
  const backendClient = new BackendClient();

  const propertyContext = getContext<PropertyContext>("selectedProperty");
</script>
{propertyContext.propertyId}
<DevlceList backendClient={backendClient} propertyId={propertyContext.propertyId} />