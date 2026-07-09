<script lang="ts">
  import type { IEntryDeviceModel } from '$lib/models/item-models';
  import { Button, TextInput, Tile } from 'carbon-components-svelte';
  import type { Snippet } from 'svelte';

  let {
    deleteButton,
    entryDevice,
    onSave,
  } = $props<{
    deleteButton?: Snippet<[() => void]>;
    entryDevice: IEntryDeviceModel;
    onSave: (ed: IEntryDeviceModel) => Promise<void>;
  }>();
</script>

<TextInput labelText="Device Type" bind:value={entryDevice.deviceType} />
<br />
<TextInput labelText="Identifier" bind:value={entryDevice.identifier} />
<br />
<TextInput labelText="Defunct Reason" bind:value={entryDevice.defunctReason} />
<br />
<br />
<Button onclick={() => onSave(entryDevice)}>Save Entry Device</Button>
{#if deleteButton !== undefined}
  {@render deleteButton(() => {})}
{/if}
