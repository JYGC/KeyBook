<svelte:options tag="app-svelte" />
<script lang="ts">
	export let name;
	export let id;
</script>

<main>
	<h1 id="{id}">Hello {name}!</h1>
</main>

<style>
    h1 {
        font-size: 5em;
    }
</style>
