<script lang="ts">
  import "carbon-components-svelte/css/all.css";

  const { children } = $props();

  let theme = "g80"; // "white" | "g10" | "g80" | "g90" | "g100"

  $effect(() => {
    document.documentElement.setAttribute("theme", theme);
  });
</script>

{@render children()}