<script lang="ts">
  import type { ICobrandEditorModule } from '$lib/modules/interfaces';
  import { Button, TextInput, Tile } from 'carbon-components-svelte';
  import type { Snippet } from 'svelte';

  let {
    deleteButton,
    cobrandEditorModule,
  } = $props<{
    deleteButton?: Snippet<[() => void]>;
    cobrandEditorModule: ICobrandEditorModule;
  }>();

  const saveButtonClick = cobrandEditorModule.getSaveCobrandAction();
  const deleteActionButtonClick = cobrandEditorModule.getDeleteCobrandAction();
</script>

{#await cobrandEditorModule.cobrandAsync}
  <Tile>...getting cobrand details</Tile>
{:then cobrand}
  {#if cobrand === null}
    <p>Cobrand not found.</p>
  {:else}
    <TextInput labelText="Cobrand Name" bind:value={cobrand.name} />
    <br />
    <br />
    <Button onclick={() => saveButtonClick(cobrand)}>Save</Button>
    {#if deleteButton !== undefined && deleteActionButtonClick !== undefined}
      {@render deleteButton(() => deleteActionButtonClick(cobrand.id))}
    {/if}
  {/if}
{:catch error}
  {error}
{/await}
