<script lang="ts">
	import type { IPersonEditorModule } from "$lib/modules/interfaces";
  import { Button, ClickableTile, Select, SelectItem, TextInput, Tile } from "carbon-components-svelte";
	import type { Snippet } from "svelte";

  let {
    children,
    deleteButton,
    personEditorModule = $bindable(),
  } = $props<{
    children?: Snippet,
    deleteButton?: Snippet<[() => null]>,
    personEditorModule: IPersonEditorModule
  }>();

  let allowEditingPersonType = $state(personEditorModule.isAdd)
  const setAllowEditingPersonType = () => allowEditingPersonType = true;

  const saveButtonClick = personEditorModule.getSavePersonAction();
  const deleteActionButtonClick = personEditorModule.getDeletePersonAction();
</script>
{#await personEditorModule.personAsync}
  <Tile>...getting person details</Tile>
{:then person}
  {#if person === null}
    {personEditorModule.callBackAction()}
  {:else}
    <TextInput labelText="Person Name" bind:value={person.name} />
    <br />
    {#if allowEditingPersonType}
      <Select labelText="Person Type" bind:selected={person.type}>
        <SelectItem value="Tenant" />
        <SelectItem value="Agent" />
        <SelectItem value="Household" />
        <SelectItem value="Owner" />
      </Select>
    {:else}
      <ClickableTile onclick={setAllowEditingPersonType}>
        Person Type: {person.type} - Change
      </ClickableTile>
    {/if}
    <br />
    {#if children !== undefined}
      {@render children()}
      <br />
    {/if}
    <br />
    <Button onclick={() => saveButtonClick(person)}>Save</Button>
    {#if deleteButton !== undefined && deleteActionButtonClick !== null}
      {@render deleteButton(() => deleteActionButtonClick(person))}
    {/if}
  {/if}
{/await}