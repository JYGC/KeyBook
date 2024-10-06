<script lang="ts">
	import { goto } from "$app/navigation";
	import type { IPersonListModule } from "$lib/interfaces";
	import { Button, DataTable, Tile } from "carbon-components-svelte";
  
  let {
    personListModule,
    propertyId,
    selectedPersonId = $bindable(),
  } = $props<{
    personListModule: IPersonListModule,
    propertyId: string,
    selectedPersonId: string,
  }>();

  const gotoPersonDetails = (personId: string) => {
    selectedPersonId = personId;
    goto("/persons/edit");
  };
</script>

{#await personListModule.personListAsync}
<Tile>...getting persons</Tile>
{:then personList}
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
{:catch error}
  {error}
{/await}