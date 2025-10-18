<script lang="ts">
	import { goto } from "$app/navigation";
	import { getBackendClient } from "$lib/api/backend-client";
  import PropertyList from "$lib/components/property/PropertyList.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import { PropertyListModule } from "$lib/modules/property/property-list-module.svelte";
	import { Button } from "carbon-components-svelte";

  const backendClient = getBackendClient();

  const propertyListModule = new PropertyListModule(backendClient);

  const selectedProperty = getPropertyContext();

  const gotoAddProperty = () => {
    goto("/user/properties/add");
  };
</script>

<Button onclick={gotoAddProperty}>Add Property</Button>
<PropertyList
propertyListModule={propertyListModule}
  bind:selectedPropertyId={selectedProperty.selectedPropertyId}
/>