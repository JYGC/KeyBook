<script lang="ts">
	import { setDeviceContext } from "$lib/contexts/device-context.svelte";

  const { children } = $props();

  setDeviceContext();
</script>

{@render children()}