<script lang="ts">
	import { BackendClient } from "$lib/api/backend-client.svelte";
	import PersonList from "$lib/components/person/PersonList.svelte";
	import { getContext } from "svelte";
	import type { PropertyContext } from "$lib/stores/property-context.svelte";

	const propertyContext = getContext<PropertyContext>("selectedProperty");

  const backendClient = new BackendClient();
</script>
{propertyContext.propertyId}
<PersonList backendClient={backendClient} propertyId={propertyContext.propertyId} />