<script lang="ts">
	import { DeviceListDataService } from "$lib/modules/device-list-data-service.svelte";

  const deviceListDataService = new DeviceListDataService();
</script>
<button onclick={deviceListDataService.add}>Add</button>

{#each deviceListDataService.devices as device}
  <div>{device}</div>
{/each}