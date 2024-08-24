<script lang="ts">
	import type { IEditPersonDto } from "$lib/dtos/person-dtos";
  import { Button, ClickableTile, Select, SelectItem, TextInput } from "carbon-components-svelte";

  let { 
    person,
    isAdd,
    savePersonAction,
    deletePersonAction = undefined,
  } = $props<{
    person: IEditPersonDto,
    isAdd: boolean,
    savePersonAction: (person: IEditPersonDto) => void,
    deletePersonAction?: (person: IEditPersonDto) => void,
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