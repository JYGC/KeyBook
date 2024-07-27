<script lang="ts">
	import { goto } from '$app/navigation';
	import { BackendClient } from '$lib/api/backend-client.svelte';
	import { setPropertyContext } from '$lib/contexts/property-context.svelte';
	import { Header, HeaderGlobalAction, HeaderUtilities } from 'carbon-components-svelte';
	import { Logout } from 'carbon-icons-svelte';

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
<Header company="SnailLabs" platformName="KeyBook">
  <HeaderUtilities>
    <HeaderGlobalAction
      iconDescription="Log out"
      tooltipAlignment="end"
      icon={Logout}
      onclick={logoutAndRedirect}
    />
    <!-- One click does not work -->
  </HeaderUtilities>
</Header>

{@render children()}