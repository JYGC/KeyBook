<script lang="ts">
	import { goto } from "$app/navigation";
	import type { IPropertyListItemDto } from "$lib/dtos/property-dtos";
	import type { IBackendClient } from "$lib/interfaces";
	import type { PropertyContext } from "$lib/stores/property-context.svelte";
	import { Button, DataTable } from "carbon-components-svelte";
	import { getContext } from "svelte";

  let {
    backendClient
  } = $props<{
    backendClient: IBackendClient
  }>();

  let propertyListAsync = $derived.by<Promise<IPropertyListItemDto[]>>(async () => {
    try {
      const response = await backendClient.pb.collection("properties").getList(1, 50, {
        filter: `owners.id ?~ "${backendClient.loggedInUser.id}"`,
        fields: "id,address"
      });
      return response.items;
    } catch (ex) {
      alert(ex);
      return [];
    }
  });

  const selectedProperty = getContext<PropertyContext>("selectedProperty");
  const goToDevicesOfProperty = (propertyId: string) => {
    selectedProperty.propertyId = propertyId;
    goto("/devices/listinproperty");
  };
  const goToPersonsOfProperty = (propertyId: string) => {
    selectedProperty.propertyId = propertyId;
    goto("/persons/listinproperty");
  };
</script>

{#await propertyListAsync}
  ...getting properties
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
        <Button onclick={() => goToDevicesOfProperty(cell.value)}>View Devices</Button>
        <Button onclick={() => goToPersonsOfProperty(cell.value)}>View Persons</Button>
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}