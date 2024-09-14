<script lang="ts">
	import { goto } from "$app/navigation";
	import type { IDeviceListItemModel } from "$lib/models/device-models";
	import { Button, DataTable } from "carbon-components-svelte";
  
  let {
    deviceList,
    propertyId,
    selectedDeviceId = $bindable(),
  } = $props<{
    deviceList: IDeviceListItemModel[],
    propertyId: string,
    selectedDeviceId: string,
  }>();
  
  const gotoDeviceDetails = (deviceId: string) => {
    selectedDeviceId = deviceId;
    goto("/devices/edit");
  };
</script>
<DataTable
  headers={[
    { key: "name", value: "Device Name"},
    { key: "identifier", value: "Indentification"},
    { key: "type", value: "Type"},
    { key: "id", empty: true },
  ]}
  rows={deviceList}
>
  <strong slot="title">Devices for PropertyId: {propertyId} <!--TODO: Get property name--></strong>
  <svelte:fragment slot="cell" let:cell>
    {#if cell.key === "id"}
      <Button onclick={() => gotoDeviceDetails(cell.value)}>Device Details</Button>
    {:else}
      {cell.value}
    {/if}
  </svelte:fragment>
</DataTable>