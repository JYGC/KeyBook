<script lang="ts">
  import type { INewPersonEditorModule } from '$lib/modules/interfaces';
  import { Button, TextInput, Tile } from 'carbon-components-svelte';
  import type { Snippet } from 'svelte';

  let {
    deleteButton,
    personEditorModule,
  } = $props<{
    deleteButton?: Snippet<[() => void]>;
    personEditorModule: INewPersonEditorModule;
  }>();

  const saveAction = personEditorModule.getSavePersonAction();
  const deleteActionButtonClick = personEditorModule.getDeletePersonAction();
</script>

{#await personEditorModule.personAsync}
  <Tile>...getting person details</Tile>
{:then person}
  {#if person === null}
    <p>Person not found.</p>
  {:else}
    <TextInput labelText="Name" bind:value={person.name} />
    <br />
    <TextInput labelText="Date of Birth" bind:value={person.DOB} />
    <br />
    <br />
    {#if saveAction !== null}
      <Button onclick={() => saveAction(person)}>Save</Button>
    {/if}
    {#if deleteButton !== undefined && deleteActionButtonClick !== null}
      {@render deleteButton(() => deleteActionButtonClick(person.id))}
    {/if}
  {/if}
{:catch error}
  {error}
{/await}
