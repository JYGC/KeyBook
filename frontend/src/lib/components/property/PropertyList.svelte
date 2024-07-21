<script lang="ts">
	import type { IPropertyListItemDto } from "$lib/dtos/property-dtos";
	import type { IBackendClient } from "$lib/interfaces";
	import { Button, DataTable, OverflowMenu, OverflowMenuItem } from "carbon-components-svelte";

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
        <Button href={`/devices/listinproperty/${cell.value}`}>View Devices</Button>
        <Button href={`/persons/listinproperty/${cell.value}`}>View Persons</Button>
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}