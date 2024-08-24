<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import DeviceEditor from "$lib/components/device/DeviceEditor.svelte";
	import { getDeviceContext } from "$lib/contexts/device-context.svelte";
	import type { IEditDeviceDto } from "$lib/dtos/device-dtos";
	import { Button, Tile } from 'carbon-components-svelte';

  const deviceContext = getDeviceContext();

  const backendClient = new BackendClient();

  let deviceAsync = $derived.by<Promise<IEditDeviceDto | null>>(async () => {
    try {
      return await backendClient.pb.collection("devices").getOne<IEditDeviceDto>(
        deviceContext.selectedDeviceId,
        { fields: "id,type,name,identifier,defunctreason,property" },
      );
    } catch (ex) {
      alert(ex);
      return null;
    }
  });

  const gotoPropertyDeviceList = () => {
    goto("/devices/list/property");
  };

  const saveDeviceActionAsync = async (changedDevice: IEditDeviceDto) => {
    try {
      await backendClient.pb.collection("devices").update(changedDevice.id, {
        type: changedDevice.type,
        name: changedDevice.name,
        identifier: changedDevice.identifier,
        defunctreason: changedDevice.defunctreason,
      });
      gotoPropertyDeviceList();
    } catch (ex) {
      alert(ex);
    }
  };

  const deleteDeviceActionAsync = async (device: IEditDeviceDto) => {
    try {
      await backendClient.pb.collection("devices").delete(device.id);
      gotoPropertyDeviceList();
    } catch (ex) {
      alert(ex);
    }
  };
</script>

<Button onclick={gotoPropertyDeviceList}>Back</Button>

{#await deviceAsync}
  <Tile>...getting device details</Tile>
{:then device}
  {#if device === null}
    {history.back()}
  {:else}
    <DeviceEditor
      device={device}
      isAdd={false}
      saveDeviceAction={saveDeviceActionAsync}
      deleteDeviceAction={deleteDeviceActionAsync}
    />
  {/if}
{:catch error}
  {error}
{/await}