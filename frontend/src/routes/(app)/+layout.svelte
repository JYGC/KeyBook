<script lang="ts">
	import { goto } from '$app/navigation';
	import { BackendClient } from '$lib/api/backend-client.svelte';
	import { PropertyContext } from '$lib/stores/property-context.svelte';
	import { setContext } from "svelte";

  const { children } = $props();

  const propertyContext = new PropertyContext();
  setContext<PropertyContext>("selectedProperty", propertyContext);

  const logoutAndRedirect = async () => {
    const authManager = new BackendClient();
    await authManager.logoutAsync();
    if (!authManager.isTokenValid) {
      goto("/auth");
    }
  }
</script>
<button onclick={logoutAndRedirect}>Logout</button>

{@render children()}