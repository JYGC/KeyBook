<script lang="ts">
	import { goto } from '$app/navigation';
	import { BackendClient } from '$lib/api/backend-client.svelte';
	import { setPropertyContext } from '$lib/contexts/property-context.svelte';

  const { children } = $props();

  setPropertyContext();

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