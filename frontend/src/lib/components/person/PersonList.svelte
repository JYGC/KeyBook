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
      { key: "personname", value: "Name" },
      { key: "persontype", value: "Type" },
      { key: "devices", value: "Holding Devices" },
      { key: "personid", empty: true },
    ]}
    rows={personList}
  >
    <strong slot="title">Persons for PropertyId: {propertyId } <!--Get property name--></strong>
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === "personid"}
        <Button onclick={() => gotoPersonDetails(cell.value)}>Person Details</Button>
      {:else if cell.key === "devices"}
        {JSON.stringify(cell.value)}
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}