<svelte:options tag="device-edit-svelte" />
<script lang="ts">
    export let devicejson;

    import DeviceDetails from './DeviceDetails.svelte';

    let device = JSON.parse(devicejson);
    console.log(device);
</script>
<main>
    <device-details-svelte {...device} disabletype=true></device-details-svelte>
</main>