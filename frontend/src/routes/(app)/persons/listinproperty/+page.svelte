<script lang="ts">
	import { BackendClient } from "$lib/api/backend-client.svelte";
	import PersonList from "$lib/components/person/PersonList.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";

	const propertyContext = getPropertyContext();

  const backendClient = new BackendClient();
</script>
{propertyContext.selectedPropertyId}
<PersonList backendClient={backendClient} propertyId={propertyContext.selectedPropertyId} />