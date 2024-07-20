<script lang="ts">
	import type { IDeviceListItem } from "$lib/dtos/device-dtos";
	import { type IBackendClient } from "../../interfaces";
  
  let {
    backendClient,
    propertyId,
  } = $props<{
    backendClient: IBackendClient,
    propertyId: string,
  }>();

  let deviceListAsync = $derived.by<Promise<IDeviceListItem[]>>(async () => {
    try {
      return await backendClient.pb.collection("devices").getList(1, 50, {
        filter: `property.id = "${propertyId}"`
      });
    } catch (ex) {
      alert(ex);
      return [];
    }
  });
</script>

{#await deviceListAsync}
  ...getting devices
{:then deviceList}
  {JSON.stringify(deviceList)}
{:catch error}
  {error}
{/await}