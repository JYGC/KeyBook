<script lang="ts">
	import type { IEditDeviceDto } from "$lib/dtos/device-dtos";
	import { Button, ClickableTile, Select, SelectItem, TextInput } from "carbon-components-svelte";

  let { 
    device,
    isAdd,
    saveDeviceAction,
    deleteDeleteAction = undefined,
  } = $props<{
    device: IEditDeviceDto,
    isAdd: boolean,
    saveDeviceAction: (device :IEditDeviceDto) => void,
    deleteDeleteAction?: (device :IEditDeviceDto) => void,
  }>()

  let showDefunctReason = $state(isAdd || false);
  const setShowDefunctReason = () => showDefunctReason = true;

  let deviceStatusText = $derived(device.defunctreason === "None" ? "Usuable" : device.defunctreason);

  const saveButtonClick = () => saveDeviceAction(device);
  const deleteButtonClick = () => deleteDeleteAction(device);
</script>
<TextInput labelText="Device Name" bind:value={device.name} />
<br />
<TextInput labelText="Device Identifier" bind:value={device.identifier} />
<br />
{#if showDefunctReason}
  <Select labelText="Device Defunct Reason" bind:selected={device.defunctreason}>
    <SelectItem value="None" />
    <SelectItem value="Lost" />
    <SelectItem value="Damaged" />
    <SelectItem value="Retired" />
    <SelectItem value="Stolen" />
  </Select>
{:else}
  <ClickableTile onclick={setShowDefunctReason}>
    Device status: {deviceStatusText} - Change
  </ClickableTile>
{/if}
<br />
<Select labelText="Device Type" bind:selected={device.type} disabled={!isAdd}>
  <SelectItem value="Fob" />
  <SelectItem value="Key" />
  <SelectItem value="Remote" />
  <SelectItem value="RoomKey" />
  <SelectItem value="MailboxKey" />
</Select>
<br />
<br />
<Button onclick={saveButtonClick}>Save</Button>
{#if deleteDeleteAction !== undefined}
  <Button onclick={deleteButtonClick}>Delete</Button>
{/if}