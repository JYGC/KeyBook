<script lang="ts">
  import DevlceList from '$lib/components/device/DevlceList.svelte';
  import { BackendClient } from '$lib/api/backend-client';
	import { getPropertyContext } from '$lib/contexts/property-context.svelte';
	import type { IDeviceListItemDto } from '$lib/dtos/device-dtos';
	import { getDeviceContext } from '$lib/contexts/device-context.svelte';
	import { goto } from '$app/navigation';
	import { Button, Tile } from 'carbon-components-svelte';

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

  const gotoPropertyList = () => {
    goto("/properties/list");
  };

const gotoAddNewDevice = () => {
  goto("/devices/add");
};
</script>

<Button onclick={gotoPropertyList}>Back</Button>
<Button onclick={gotoAddNewDevice}>Add New Device</Button>

{#await deviceListAsync}
<Tile>...getting devices</Tile>
{:then deviceList}
  <DevlceList
    deviceList={deviceList}
    propertyId={propertyContext.selectedPropertyId}
    bind:selectedDeviceId={deviceContext.selectedDeviceId}
  />
{:catch error}
  {error}
{/await}