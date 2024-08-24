<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
  import PropertyList from "$lib/components/property/PropertyList.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import type { IPropertyListItemDto } from "$lib/dtos/property-dtos";
	import { Button, OverflowMenu, OverflowMenuItem, Tile } from "carbon-components-svelte";

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

  const gotoImportCsvDate = () => {
    goto("/importdata/csv1");
  };

  const gotoAddProperty = () => {
    goto("/properties/add");
  };
</script>

<OverflowMenu size="xl" style="width: auto;">
  <div slot="menu">
    <Button onclick={gotoAddProperty}>Add Property</Button>
  </div>
  <OverflowMenuItem text="Add One Property" />
  <OverflowMenuItem onclick={gotoImportCsvDate} text="Add via CSV" />
</OverflowMenu>
{#await propertyListAsync}
<Tile>...getting properties</Tile>
{:then propertyList}
  <PropertyList
    propertyList={propertyList}
    bind:selectedPropertyId={selectedProperty.selectedPropertyId}
  />
{:catch error}
  {error}
{/await}