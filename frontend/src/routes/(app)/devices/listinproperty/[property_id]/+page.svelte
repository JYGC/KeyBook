<script lang="ts">
  import DevlceList from '$lib/components/device/DevlceList.svelte';
  import type { PageData } from './$types';
  import { BackendClient } from '$lib/api/backend-client.svelte';
  
  let { data } = $props<{ data: PageData }>();
  
  const backendClient = new BackendClient();
</script>

<DevlceList backendClient={backendClient} propertyId={data.propertyId} />