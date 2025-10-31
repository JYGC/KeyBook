<script lang="ts">
	import { goto } from "$app/navigation";
	import type { IDeviceHoldingListModule } from "$lib/modules/interfaces";
	import { Button, DataTable } from "carbon-components-svelte";

  let {
    deviceHoldingListModule,
    selectedDeviceId = $bindable(),
  } = $props<{
    deviceHoldingListModule: IDeviceHoldingListModule,
    selectedDeviceId: string,
  }>();

  const gotoDeviceDetails = (deviceId: string) => {
    selectedDeviceId = deviceId;
    goto("/user/devices/edit");
  }
</script>

{#await deviceHoldingListModule.deviceHoldingListOfPersonAsync then deviceHoldingListOfPerson}
  {#if deviceHoldingListOfPerson !== null && deviceHoldingListOfPerson.length > 0}
    <DataTable
      headers={[
        { key: "name", value: "Device Name"},
        { key: "type", value: "Type"},
        { key: "id", empty: true },
      ]}
      rows={deviceHoldingListOfPerson}
    >
      <svelte:fragment slot="cell" let:cell>
        {#if cell.key === "id"}
          <Button onclick={() => gotoDeviceDetails(cell.value)}>View device</Button>
        {:else}
          {cell.value}
        {/if}
      </svelte:fragment>
    </DataTable>
  {/if}
{/await}