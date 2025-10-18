<script lang="ts">
	import PocketBase from 'pocketbase';
	import { goto } from '$app/navigation';
	import { setPropertyContext } from '$lib/contexts/property-context.svelte';
	import { setPersonContext } from "$lib/contexts/person-context.svelte";
	import { setDeviceContext } from '$lib/contexts/device-context.svelte';
	import { Button, Content, Header, HeaderUtilities } from 'carbon-components-svelte';
	import { Logout } from 'carbon-icons-svelte';
	import type { Snippet } from 'svelte';

  const {
    data,
    children
  }: {
    data: { backendClient: PocketBase },
    children: Snippet<[]>
  } = $props();

  setPropertyContext();
  setPersonContext();
  setDeviceContext();

  const backendClient = data.backendClient;

  const logoutAndRedirect = async () => {
    await backendClient.authStore.clear();
    if (!backendClient.authStore.isValid) {
      goto("/auth");
    }
  }
</script>
<Header company="SnailLabs" platformName="KeyBook">
  <HeaderUtilities>
    <Button
      iconDescription="Log out"
      tooltipAlignment="end"
      icon={Logout}
      kind="ghost"
      onclick={logoutAndRedirect}
    />
    <!-- One click does not work -->
  </HeaderUtilities>
</Header>

<Content>
  {@render children()}
</Content>