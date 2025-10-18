<script lang="ts">
	import type { IPropertyEditorModule } from "$lib/modules/interfaces";
	import { Button, TextInput, Tile } from "carbon-components-svelte";
	import type { Snippet } from "svelte";

  let { 
    deleteButton,
    propertyEditorModule
  } = $props<{
    deleteButton?: Snippet<[() => null]>,
    propertyEditorModule: IPropertyEditorModule,
  }>();

const saveButtonClick = propertyEditorModule.getSavePropertyAction();
const deleteActionButtonClick = propertyEditorModule.getDeletePropertyAction();
</script>

{#await propertyEditorModule.propertyAsync}
  <Tile>...getting property details</Tile>
{:then property}
  {#if property === null}
    {history.back()}
  {:else}
    <TextInput labelText="Property Address" bind:value={property.address} />
    <br />
    <br />
    <Button onclick={() => saveButtonClick(property)}>Save</Button>
    {#if deleteButton !== undefined && deleteActionButtonClick !== undefined}
      {@render deleteButton(() => deleteActionButtonClick(property))}
    {/if}
  {/if}
{:catch error}
  {error}
{/await}