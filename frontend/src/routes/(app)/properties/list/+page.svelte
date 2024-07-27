<script lang="ts">
	import { BackendClient } from "$lib/api/backend-client.svelte";
  import PropertyList from "$lib/components/property/PropertyList.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import type { IPropertyListItemDto } from "$lib/dtos/property-dtos";

  const backendClient = new BackendClient();

  let propertyListAsync = $derived.by<Promise<IPropertyListItemDto[]>>(async () => {
    try {
      if (backendClient.loggedInUser === null) {
        throw new Error("Cannot find loggedInUser.");
      }
      const response = await backendClient.pb.collection("properties").getList<IPropertyListItemDto>(1, 50, {
        filter: `owners.id ?~ "${backendClient.loggedInUser.id}"`,
        fields: "id,address"
      });
      return response.items;
    } catch (ex) {
      alert(ex);
      return [];
    }
  });

  const selectedProperty = getPropertyContext();
</script>


{#await propertyListAsync}
  ...getting properties
{:then propertyList}
  <PropertyList
    propertyList={propertyList}
    bind:selectedPropertyId={selectedProperty.selectedPropertyId}
  />
{:catch error}
  {error}
{/await}

<a href="/importdata/csv1">import csv</a>