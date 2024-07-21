<script lang="ts">
	import { BackendClient } from "$lib/api/backend-client.svelte";
  import PropertyList from "$lib/components/property/PropertyList.svelte";

  const backendClient = new BackendClient();
</script>

<PropertyList backendClient={backendClient} />

<a href="/importdata/csv1">import csv</a>