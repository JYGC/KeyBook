<script lang="ts">
	import type { IDeviceHistoryListModule } from "$lib/interfaces";
	import { DataTable } from "carbon-components-svelte";

  let { 
    deviceHistoryListModule,
  } = $props<{
    deviceHistoryListModule: IDeviceHistoryListModule
  }>();
</script>

{#await deviceHistoryListModule.deviceHistoryListAsync then deviceHistoryList}
  {#if deviceHistoryList !== null && deviceHistoryList.length > 0}
    <DataTable
      headers={[
        { key: "created", value: "Date" },
        { key: "description", value: "Description" },
      ]}
      rows={deviceHistoryList}
    />
  {/if}
{/await}