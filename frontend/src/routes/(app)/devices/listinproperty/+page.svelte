<script lang="ts">
  import DevlceList from '$lib/components/device/DevlceList.svelte';
  import { BackendClient } from '$lib/api/backend-client.svelte';
	import { getPropertyContext } from '$lib/contexts/property-context.svelte';
	import type { IDeviceListItemDto } from '$lib/dtos/device-dtos';
	import { getDeviceContext } from '$lib/contexts/device-context.svelte';
	import { goto } from '$app/navigation';

  const propertyContext = getPropertyContext();
  if (
    propertyContext.selectedPropertyId === null ||
    propertyContext.selectedPropertyId.trim() === ""
  ) {
    goto("/");
  }

  const deviceContext = getDeviceContext();
  
  const backendClient = new BackendClient();

  let deviceListAsync = $derived.by<Promise<IDeviceListItemDto[]>>(async () => {
    try {
      const response = await backendClient.pb.collection("devices").getList<IDeviceListItemDto>(1, 50, {
        filter: `property.id = "${propertyContext.selectedPropertyId}"`,
        fields: "id,type,name,identifier",
      });
      return response.items;
    } catch (ex) {
      alert(ex);
      return [];
    }
  });

</script>

{#await deviceListAsync}
  ...getting devices
{:then deviceList}
  <DevlceList
    deviceList={deviceList}
    propertyId={propertyContext.selectedPropertyId} 
    bind:selectedDeviceId={deviceContext.selectedDeviceId}
  />
{:catch error}
  {error}
{/await}