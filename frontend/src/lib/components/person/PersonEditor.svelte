<script lang="ts">
	import type { IPersonEditorModule } from "$lib/interfaces";
  import { Button, ClickableTile, Select, SelectItem, TextInput, Tile } from "carbon-components-svelte";

  let {
    personEditorModule = $bindable(),
  } = $props<{
    personEditorModule: IPersonEditorModule
  }>();

  let allowEditingPersonType = $state(personEditorModule.isAdd)
  const setAllowEditingPersonType = () => allowEditingPersonType = true;

  const saveButtonClick = personEditorModule.getSavePersonAction();
  const deleteButtonClick = personEditorModule.getDeletePersonAction();
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
    <br />
    <Button onclick={() => saveButtonClick(person)}>Save</Button>
    {#if deleteButtonClick !== null}
      <Button onclick={() => deleteButtonClick(person)}>Delete</Button>
    {/if}
  {/if}
{/await}