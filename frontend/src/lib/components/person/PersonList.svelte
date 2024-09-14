<script lang="ts">
	import { goto } from "$app/navigation";
	import type { IPersonIdNameTypeModel } from "$lib/models/person-models";
	import { Button, DataTable } from "carbon-components-svelte";
  
  let {
    personList,
    propertyId,
    selectedPersonId = $bindable(),
  } = $props<{
    personList: IPersonIdNameTypeModel[],
    propertyId: string,
    selectedPersonId: string,
  }>();

  const gotoPersonDetails = (personId: string) => {
    selectedPersonId = personId;
    goto("/persons/edit");
  };
</script>
<DataTable
  headers={[
    { key: "name", value: "Name" },
    { key: "type", value: "Type" },
    { key: "id", empty: true },
  ]}
  rows={personList}
>
  <strong slot="title">Persons for PropertyId: {propertyId } <!--Get property name--></strong>
  <svelte:fragment slot="cell" let:cell>
    {#if cell.key === "id"}
      <Button onclick={() => gotoPersonDetails(cell.value)}>Person Details</Button>
    {:else}
      {cell.value}
    {/if}
  </svelte:fragment>
</DataTable>