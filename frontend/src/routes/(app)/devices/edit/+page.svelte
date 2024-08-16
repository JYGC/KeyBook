<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import DeviceEdit from "$lib/components/device/DeviceEdit.svelte";
	import { getDeviceContext } from "$lib/contexts/device-context.svelte";
	import type { IDeviceEditDto } from "$lib/dtos/device-dtos";
	import { Button, Tile } from 'carbon-components-svelte';

  const deviceContext = getDeviceContext();

  const backendClient = new BackendClient();

  let deviceAsync = $derived.by<Promise<IDeviceEditDto | null>>(async () => {
    try {
      return await backendClient.pb.collection("devices").getOne<IDeviceEditDto>(
        deviceContext.selectedDeviceId,
        { fields: "id,type,name,identifier,defunctreason,property" }
      )
    } catch (ex) {
      alert(ex);
      return null;
    }
  });

  const saveDeviceActionAsync = async (changedDevice: IDeviceEditDto) => {
    try {
      await backendClient.pb.collection("devices").update(changedDevice.id, {
        type: changedDevice.type,
        name: changedDevice.name,
        identifier: changedDevice.identifier,
        defunctreason: changedDevice.defunctreason,
      });
      goto("/devices/listinproperty");
    } catch (ex) {
      alert(ex)
    }
  };

  const gotoPropertyList = () => {
    goto("/devices/listinproperty");
  };
</script>

<Button onclick={gotoPropertyList}>Back</Button>

{#await deviceAsync}
  <Tile>...getting device details</Tile>
{:then device}
  {#if device === null}
    {history.back()}
  {:else}
    <DeviceEdit
      device={device}
      isAdd={false}
      saveDeviceAction={saveDeviceActionAsync}
    />
  {/if}
{:catch error}
  {error}
{/await}