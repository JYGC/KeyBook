<script lang="ts">
	import type { IDeviceEditorModule } from "$lib/interfaces";
	import { Button, ClickableTile, Select, SelectItem, TextInput, Tile } from "carbon-components-svelte";
	import type { Snippet } from "svelte";

  let {
    children,
    deleteButton,
    deviceEditorModule = $bindable(),
  } = $props<{
    children?: Snippet,
    deleteButton?: Snippet<[() => null]>,
    deviceEditorModule: IDeviceEditorModule,
  }>();

  let showDefunctReason = $state(deviceEditorModule.isAdd);
  const setShowDefunctReason = () => showDefunctReason = true;

  const saveButtonClick = deviceEditorModule.getSaveDeviceAction();
  const deleteActionButtonClick = deviceEditorModule.getDeleteDeviceAction();
</script>
{#await deviceEditorModule.deviceAsync}
  <Tile>...getting device details</Tile>
{:then device}
  {#if device === null}
    {deviceEditorModule.callBackAction()}
  {:else}
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
    {#await deviceEditorModule.deviceStatusTextAsync then deviceStatusText}
      <ClickableTile onclick={setShowDefunctReason}>
        Device status: {deviceStatusText} - Change
      </ClickableTile>
    {/await}
    {/if}
    <br />
    {#if children !== undefined}
      {@render children()}
      <br />
    {/if}
    <Select labelText="Device Type" bind:selected={device.type} disabled={!deviceEditorModule.isAdd}>
      <SelectItem value="Key" />
      <SelectItem value="RoomKey" />
      <SelectItem value="MailboxKey" />
      <SelectItem value="Fob" />
      <SelectItem value="Remote" />
    </Select>
    <br />
    <br />
    <Button onclick={() => saveButtonClick(device)}>Save</Button>
    {#if deleteButton !== undefined && deleteActionButtonClick !== undefined}
      {@render deleteButton(() => deleteActionButtonClick(device))}
    {/if}
  {/if}
{/await}