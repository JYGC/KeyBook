<script lang="ts">
  import { getBackendClient } from '$lib/api/backend-client';
  import ItemEditor from '$lib/components/item/ItemEditor.svelte';
  import EntryDeviceEditor from '$lib/components/item/EntryDeviceEditor.svelte';
  import ConfirmButtonAndDialog from '$lib/components/shared/ConfirmButtonAndDialog.svelte';
  import { getItemContext } from '$lib/contexts/item-context.svelte';
  import { ItemRepository } from '$lib/repositories/item/item-repository';
  import { EntryDeviceRepository } from '$lib/repositories/item/entry-device-repository';
  import { ItemService } from '$lib/services/item/item-service';
  import { ItemDetailModule } from '$lib/modules/item/item-detail-module.svelte';
  import { Button, TextInput, Tile } from 'carbon-components-svelte';

  const itemContext = getItemContext();

  const goBack = () => {
    history.back();
  };

  const pb = getBackendClient();
  const itemRepo = new ItemRepository(pb);
  const entryDeviceRepo = new EntryDeviceRepository(pb);
  const itemService = new ItemService(itemRepo, entryDeviceRepo);
  const itemDetailModule = new ItemDetailModule(itemService, itemContext, goBack);

  const saveEntryDevice = itemDetailModule.getSaveEntryDeviceAction();
  const deleteEntryDevice = itemDetailModule.getDeleteEntryDeviceAction();

  let newDeviceType = $state('');
  let newIdentifier = $state('');
  let newDefunctReason = $state('');

  const createEntryDevice = itemDetailModule.getCreateEntryDeviceAction();
</script>

<Button onclick={goBack}>Back</Button>

<ItemEditor itemEditorModule={itemDetailModule}>
  {#snippet deleteButton(deleteActionButtonClick: () => void)}
    <ConfirmButtonAndDialog
      submitAction={deleteActionButtonClick}
      buttonText="Delete Item"
      bodyMessage="Are you sure you want to delete this item? This will also delete its entry device."
    />
  {/snippet}
</ItemEditor>

<br />

{#await itemDetailModule.entryDeviceAsync}
  <Tile>...getting entry device</Tile>
{:then entryDevice}
  {#if entryDevice !== null}
    <h3>Entry Device</h3>
    <EntryDeviceEditor
      {entryDevice}
      onSave={saveEntryDevice}
    >
      {#snippet deleteButton(_: () => void)}
        <ConfirmButtonAndDialog
          submitAction={() => deleteEntryDevice(entryDevice.id)}
          buttonText="Delete Entry Device"
          bodyMessage="Are you sure you want to remove the entry device from this item?"
        />
      {/snippet}
    </EntryDeviceEditor>
  {:else}
    <h3>Add Entry Device</h3>
    <TextInput labelText="Device Type" bind:value={newDeviceType} />
    <br />
    <TextInput labelText="Identifier" bind:value={newIdentifier} />
    <br />
    <TextInput labelText="Defunct Reason" bind:value={newDefunctReason} />
    <br />
    <br />
    <Button onclick={() => createEntryDevice(itemContext.selectedItemId, newDeviceType, newIdentifier, newDefunctReason)}>
      Create Entry Device
    </Button>
  {/if}
{:catch error}
  {error}
{/await}
