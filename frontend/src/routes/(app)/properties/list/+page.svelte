<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
  import PropertyList from "$lib/components/property/PropertyList.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import type { IPropertyListItemModel } from "$lib/models/property-models";
	import { PropertyListModule } from "$lib/modules/property/property-list-module.svelte";
	import { Button, OverflowMenu, OverflowMenuItem, Tile } from "carbon-components-svelte";

  const backendClient = new BackendClient();

  const propertyListModule = new PropertyListModule(backendClient);

  const selectedProperty = getPropertyContext();

  const gotoAddProperty = () => {
    goto("/properties/add");
  };
</script>

<OverflowMenu size="xl" style="width: auto;">
  <div slot="menu">
    <Button>Add Property</Button>
  </div>
  <OverflowMenuItem onclick={gotoAddProperty} text="Add One Property" />
</OverflowMenu>
<PropertyList
propertyListModule={propertyListModule}
  bind:selectedPropertyId={selectedProperty.selectedPropertyId}
/>