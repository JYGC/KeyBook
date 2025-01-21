<script lang="ts">
	import type { IEditPropertyModel } from "$lib/models/property-models";
	import { Button, TextInput } from "carbon-components-svelte";
	import type { Snippet } from "svelte";

  let { 
    deleteButton,
    property,
    isAdd,
    savePropertyAction,
    deletePropertyAction = undefined,
  } = $props<{
    deleteButton?: Snippet<[() => null]>,
    property: IEditPropertyModel,
    isAdd: boolean,
    savePropertyAction: (property: IEditPropertyModel) => void,
    deletePropertyAction?: (property: IEditPropertyModel) => void,
  }>();

const saveButtonClick = () => savePropertyAction(property);
const deleteActionButtonClick = () => deletePropertyAction(property);
</script>
<TextInput labelText="Property Address" bind:value={property.address} />
<br />
<br />
<Button onclick={saveButtonClick}>Save</Button>
{#if deleteButton !== undefined && deletePropertyAction !== undefined}
{@render deleteButton(() => deleteActionButtonClick())}
{/if}