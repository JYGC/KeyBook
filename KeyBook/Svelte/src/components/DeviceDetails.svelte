<svelte:options tag="device-details-svelte" />
<script lang="ts">
    export let name;
    export let identifier;
    export let status;
    export let type;
    export let hidestatus = false;
    export let disabletype = false;
</script>
<main>
    {name}
    {identifier}
    {status}
    {type}
    {hidestatus}
    {disabletype}
    <!--continue here-->
</main>
