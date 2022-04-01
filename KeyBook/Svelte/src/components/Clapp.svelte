<svelte:options tag="clapp-svelte" />
<script lang="ts">
    export let num1;
    export let num2;

    let c;
    let d;
    $: {
        c = Number(num1) + Number(num2);
        d = Number(num1) * Number(num2);
    }
</script>

<main>
    <h1>num1 + num2 = {c}</h1>
    <h1>num1 * num2 = {d}</h1>
</main>

<style>
    h1 {
        font-size: 5em;
    }
</style>
