<script lang="ts">
	import { DeviceListModule } from '$lib/modules/device/device-list-module.svelte';
  import DevlceList from '$lib/components/device/DevlceList.svelte';
  import { BackendClient } from '$lib/api/backend-client';
	import { getPropertyContext } from '$lib/contexts/property-context.svelte';
	import { getDeviceContext } from '$lib/contexts/device-context.svelte';
	import { goto } from '$app/navigation';
	import { Button } from 'carbon-components-svelte';

  const propertyContext = getPropertyContext();
  if (
    propertyContext.selectedPropertyId === null ||
    propertyContext.selectedPropertyId.trim() === ""
  ) {
    goto("/");
  }

  const deviceContext = getDeviceContext();

  const backendClient = new BackendClient();
  const deviceListModule = new DeviceListModule(backendClient, propertyContext);

  const gotoPropertyList = () => {
    goto("/properties/list");
  };

const gotoAddNewDevice = () => {
  goto("/devices/add");
};
</script>

<Button onclick={gotoPropertyList}>Back</Button>
<Button onclick={gotoAddNewDevice}>Add New Device</Button>

<DevlceList
  deviceListModule={deviceListModule}
  propertyId={propertyContext.selectedPropertyId}
  bind:selectedDeviceId={deviceContext.selectedDeviceId}
/>