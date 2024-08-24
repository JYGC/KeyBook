<script lang="ts">
	import type { IEditPropertyDto } from "$lib/dtos/property-dtos";
	import { Button, TextInput } from "carbon-components-svelte";

  let { 
    property,
    isAdd,
    savePropertyAction,
    deletePropertyAction = undefined,
  } = $props<{
    property: IEditPropertyDto,
    isAdd: boolean,
    savePropertyAction: (property: IEditPropertyDto) => void,
    deletePropertyAction?: (property: IEditPropertyDto) => void,
  }>();

const saveButtonClick = () => savePropertyAction(property);
const deleteButtonClick = () => deletePropertyAction(property);
</script>
<TextInput labelText="Property Address" bind:value={property.address} />
<br />
<br />
<Button onclick={saveButtonClick}>Save</Button>
{#if deletePropertyAction !== undefined}
  <Button onclick={deleteButtonClick}>Delete</Button>
{/if}