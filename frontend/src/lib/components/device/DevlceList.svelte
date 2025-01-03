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
    goto("/devices/edit");
  };
</script>

{#await deviceListModule.deviceListAsync}
  <Tile>...getting devices</Tile>
{:then deviceList}
  <DataTable
    headers={[
      { key: "devicename", value: "Device Name"},
      { key: "deviceidentifier", value: "Indentification"},
      { key: "devicetype", value: "Type"},
      { key: "personname", value: "Current Device Holder"},
      { key: "deviceid", empty: true },
    ]}
    rows={deviceList}
  >
    <strong slot="title">Devices for PropertyId: {propertyId} <!--TODO: Get property name--></strong>
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === "deviceid"}
        <Button onclick={() => gotoDeviceDetails(cell.value)}>Device Details</Button>
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}