<script lang="ts">
	import { goto } from "$app/navigation";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import { Button } from "carbon-components-svelte";

  const propertyContext = getPropertyContext();

  const propertyId = propertyContext.selectedPropertyId;

	const gotoPropertyListAsync = async () => {
    goto("/devices/listinproperty");
	}
</script>
<Button onclick={gotoPropertyListAsync}>Back</Button>