<script lang="ts">
	import { goto } from "$app/navigation";
	import type { IPersonListModule } from "$lib/modules/interfaces";
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
    goto("/user/persons/edit");
  };
</script>

{#await personListModule.personListAsync}
<Tile>...getting persons</Tile>
{:then personList}
  <DataTable
    headers={[
      { key: "personName", value: "Name" },
      { key: "personType", value: "Type" },
      { key: "holdingDeviceJsons", value: "Holding Devices Json" },
      { key: "personId", empty: true },
    ]}
    rows={personList}
  >
    <strong slot="title">Persons for PropertyId: {propertyId } <!--Get property name--></strong>
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === "personId"}
        <Button onclick={() => gotoPersonDetails(cell.value)}>Person Details</Button>
      {:else if cell.key === "holdingDeviceJsons"}
        {#if cell.value !== null && cell.value instanceof Array}
          {#each cell.value as holdingDevice}
            {#if "deviceName" in holdingDevice}
              <div>{holdingDevice.deviceName}</div>
            {/if}
          {/each}
        {/if}
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}