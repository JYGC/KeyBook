<script lang="ts">
	import type { IEditPersonModel } from "$lib/models/person-models";
  import { Button, ClickableTile, Select, SelectItem, TextInput } from "carbon-components-svelte";

  let { 
    person,
    isAdd,
    savePersonAction,
    deletePersonAction = undefined,
  } = $props<{
    person: IEditPersonModel,
    isAdd: boolean,
    savePersonAction: (person: IEditPersonModel) => void,
    deletePersonAction?: (person: IEditPersonModel) => void,
  }>();

  let allowEditingPersonType = $state(isAdd)
  const setAllowEditingPersonType = () => allowEditingPersonType = true;

  const saveButtonClick = () => savePersonAction(person);
  const deleteButtonClick = () => deletePersonAction(person);
</script>
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
<Button onclick={saveButtonClick}>Save</Button>
{#if deletePersonAction !== undefined}
  <Button onclick={deleteButtonClick}>Delete</Button>
{/if}