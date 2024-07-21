<script lang="ts">
	import { BackendClient } from "$lib/api/backend-client.svelte";
	import PersonList from "$lib/components/person/PersonList.svelte";
	import type { PageData } from "../../$types";

  let { data } = $props<{ data: PageData }>();

  const backendClient = new BackendClient();
</script>

<PersonList backendClient={backendClient} propertyId={data.propertyId} />