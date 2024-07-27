<script lang="ts">
	import { setPersonContext } from "$lib/contexts/person-context.svelte";

  const { children } = $props();

  setPersonContext();
</script>

{@render children()}