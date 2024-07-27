<script lang="ts">
	import type { IPersonListItemDto } from "$lib/dtos/person-dtos";
	import { Button, DataTable } from "carbon-components-svelte";
  
  let {
    personList,
    propertyId,
    selectedPersonId = $bindable(),
  } = $props<{
    personList: IPersonListItemDto[],
    propertyId: string,
    selectedPersonId: string,
  }>();

  const gotoPersonDetails = (personId: string) => {
    selectedPersonId = personId;
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