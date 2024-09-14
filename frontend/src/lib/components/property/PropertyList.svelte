<script lang="ts">
	import { goto } from "$app/navigation";
	import type { IPropertyListItemModel } from "$lib/models/property-models";
	import { Button, ButtonSet, DataTable } from "carbon-components-svelte";

  let {
    propertyList,
    selectedPropertyId = $bindable(),
  } = $props<{
    propertyList: IPropertyListItemModel[],
    selectedPropertyId: string,
  }>();

  const goToDevicesOfProperty = (propertyId: string) => {
    selectedPropertyId = propertyId;
    goto("/devices/list/property");
  };
  const goToPersonsOfProperty = (propertyId: string) => {
    selectedPropertyId = propertyId;
    goto("/persons/list/property");
  };
  const goToEditProperty = (propertyId: string) => {
    selectedPropertyId = propertyId;
    goto("/properties/edit");
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