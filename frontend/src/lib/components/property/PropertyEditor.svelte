<script lang="ts">
	import type { IEditPropertyModel } from "$lib/dtos/property-models";
	import { Button, TextInput } from "carbon-components-svelte";

  let { 
    property,
    isAdd,
    savePropertyAction,
    deletePropertyAction = undefined,
  } = $props<{
    property: IEditPropertyModel,
    isAdd: boolean,
    savePropertyAction: (property: IEditPropertyModel) => void,
    deletePropertyAction?: (property: IEditPropertyModel) => void,
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