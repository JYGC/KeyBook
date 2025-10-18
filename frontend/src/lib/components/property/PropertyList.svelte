<script lang="ts">
	import { goto } from "$app/navigation";
	import type { IPropertyListModule } from "$lib/interfaces";
	import { Button, ButtonSet, DataTable, Tile } from "carbon-components-svelte";

  let {
    propertyListModule,
    selectedPropertyId = $bindable(),
  } = $props<{
    propertyListModule: IPropertyListModule,
    selectedPropertyId: string,
  }>();

  const goToDevicesOfProperty = (propertyId: string) => {
    selectedPropertyId = propertyId;
    goto("/user/devices/list/property");
  };
  const goToPersonsOfProperty = (propertyId: string) => {
    selectedPropertyId = propertyId;
    goto("/user/persons/list/property");
  };
  const goToEditProperty = (propertyId: string) => {
    selectedPropertyId = propertyId;
    goto("/user/properties/edit");
  };
</script>
{#await propertyListModule.propertyListAsync}
<Tile>...getting properties</Tile>
{:then propertyList}
  <DataTable
    headers={[
      { key: "address", value: "Property Address" },
      { key: "id", empty: true },
    ]}
    rows={propertyList}
  >
    <strong slot="title">Properties</strong>
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === "id"}
        <ButtonSet>
          <Button onclick={() => goToDevicesOfProperty(cell.value)}>View Devices</Button>
          <Button onclick={() => goToPersonsOfProperty(cell.value)}>View Persons</Button>
          <Button onclick={() => goToEditProperty(cell.value)}>Edit</Button>
        </ButtonSet>
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}