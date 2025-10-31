<script lang="ts">
	import { DeviceListModule } from '$lib/modules/device/device-list-module.svelte';
	import { goto } from "$app/navigation";
	import { Button, DataTable, Tile } from "carbon-components-svelte";
  
  let {
    deviceListModule,
    propertyId,
    selectedDeviceId = $bindable(),
  } = $props<{
    deviceListModule: DeviceListModule,
    propertyId: string,
    selectedDeviceId: string,
  }>();
  
  const gotoDeviceDetails = (deviceId: string) => {
    selectedDeviceId = deviceId;
    goto("/user/devices/edit");
  };
</script>

{#await deviceListModule.deviceListAsync}
  <Tile>...getting devices</Tile>
{:then deviceList}
  <DataTable
    headers={[
      { key: "deviceName", value: "Device Name"},
      { key: "deviceIdentifier", value: "Identification"},
      { key: "deviceType", value: "Type"},
      { key: "personName", value: "Current Device Holder"},
      { key: "deviceId", empty: true },
    ]}
    rows={deviceList}
  >
    <strong slot="title">Devices for PropertyId: {propertyId} <!--TODO: Get property name--></strong>
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === "deviceId"}
        <Button onclick={() => gotoDeviceDetails(cell.value)}>Device Details</Button>
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}