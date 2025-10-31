<script lang="ts">
	import type { IDeviceHolderEditorModule } from "$lib/modules/interfaces";
	import { ClickableTile, Modal, Select, SelectItem } from "carbon-components-svelte";

  let { deviceHolderEditorModule } = $props<{ deviceHolderEditorModule: IDeviceHolderEditorModule }>();

  let isChangeDeviceHolderModalOpen = $state<boolean>(false);
  
  const onSubmitClick = () => {
    deviceHolderEditorModule.replaceDeviceHolderActionAsync();
    isChangeDeviceHolderModalOpen = false;
  };
</script>

<ClickableTile onclick={() => isChangeDeviceHolderModalOpen = true}>
  {#await deviceHolderEditorModule.currentDeviceHolderNameAsync then currentDeviceHolderName}
    {#if currentDeviceHolderName === ""}
      No one is currently holding this device
    {:else}
      Current device holder: {currentDeviceHolderName}
    {/if}
  {/await}
</ClickableTile>

<Modal
  bind:open={isChangeDeviceHolderModalOpen}
  modalHeading="Change Device Holder"
  primaryButtonText="Change Device Holder"
  secondaryButtonText="Cancel"
  on:click:button--secondary={() => isChangeDeviceHolderModalOpen = false}
  on:submit={onSubmitClick}
>
  <Select
    labelText="Current Device Holder"
    bind:selected={deviceHolderEditorModule.selectedDeviceHolderId}
  >
    <SelectItem value="" text="-- No holder --" />
    {#await deviceHolderEditorModule.availablePersonsAsync then availablePersons}
      {#each availablePersons as person}
        <SelectItem value={person.id} text={`${person.name} - ${person.type}`} />
      {/each}
    {/await}
  </Select>
</Modal>