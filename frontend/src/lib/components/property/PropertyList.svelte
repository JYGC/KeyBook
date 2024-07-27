<script lang="ts">
	import { goto } from "$app/navigation";
	import type { IPropertyListItemDto } from "$lib/dtos/property-dtos";
	import { Button, DataTable } from "carbon-components-svelte";

  let {
    propertyList,
    selectedPropertyId = $bindable(),
  } = $props<{
    propertyList: IPropertyListItemDto[],
    selectedPropertyId: string,
  }>();

  const goToDevicesOfProperty = (propertyId: string) => {
    selectedPropertyId = propertyId;
    goto("/devices/listinproperty");
  };
  const goToPersonsOfProperty = (propertyId: string) => {
    selectedPropertyId = propertyId;
    goto("/persons/listinproperty");
  };
</script>
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
        <Button onclick={() => goToDevicesOfProperty(cell.value)}>View Devices</Button>
        <Button onclick={() => goToPersonsOfProperty(cell.value)}>View Persons</Button>
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>