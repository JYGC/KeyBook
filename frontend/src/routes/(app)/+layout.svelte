<script lang="ts">
	import { BackendClient } from '$lib/api/backendclient.svelte';

  const { children } = $props();

  const logoutAndRedirect = async () => {
    const authManager = new BackendClient();
    await authManager.logoutAsync();
    if (!authManager.isTokenValid) {
      window.location.replace("/auth");
    }
  }
</script>
<button onclick={logoutAndRedirect}>Logout</button>

{@render children()}