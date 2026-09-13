<script lang="ts">
  import type { INewPropertyEditorModule } from '$lib/modules/interfaces';
  import { Button, TextInput, Tile } from 'carbon-components-svelte';
  import type { Snippet } from 'svelte';

  let {
    deleteButton,
    propertyEditorModule,
  } = $props<{
    deleteButton?: Snippet<[() => void]>;
    propertyEditorModule: INewPropertyEditorModule;
  }>();

  const saveAction = propertyEditorModule.getSavePropertyAction();
  const deleteActionButtonClick = propertyEditorModule.getDeletePropertyAction();
</script>

{#await propertyEditorModule.propertyAsync}
  <Tile>...getting property details</Tile>
{:then property}
  {#if property === null}
    <p>Property not found.</p>
  {:else}
    <TextInput labelText="Property Address" bind:value={property.address} />
    <br />
    <br />
    {#if saveAction !== null}
      <Button onclick={() => saveAction(property)}>Save</Button>
    {/if}
    {#if deleteButton !== undefined && deleteActionButtonClick !== null}
      {@render deleteButton(() => deleteActionButtonClick(property.id))}
    {/if}
  {/if}
{:catch error}
  {error}
{/await}
