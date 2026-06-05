<script lang="ts">
  import type { IItemEditorModule } from '$lib/modules/interfaces';
  import { Button, TextInput, Tile } from 'carbon-components-svelte';
  import type { Snippet } from 'svelte';

  let {
    deleteButton,
    itemEditorModule,
  } = $props<{
    deleteButton?: Snippet<[() => void]>;
    itemEditorModule: IItemEditorModule;
  }>();

  const saveButtonClick = itemEditorModule.getSaveItemAction();
  const deleteActionButtonClick = itemEditorModule.getDeleteItemAction();
</script>

{#await itemEditorModule.itemAsync}
  <Tile>...getting item details</Tile>
{:then item}
  {#if item === null}
    {history.back()}
  {:else}
    <TextInput labelText="Item Name" bind:value={item.name} />
    <br />
    <TextInput labelText="Description" bind:value={item.description} />
    <br />
    <br />
    <Button onclick={() => saveButtonClick(item)}>Save</Button>
    {#if deleteButton !== undefined && deleteActionButtonClick !== undefined}
      {@render deleteButton(() => deleteActionButtonClick(item.id))}
    {/if}
  {/if}
{:catch error}
  {error}
{/await}
