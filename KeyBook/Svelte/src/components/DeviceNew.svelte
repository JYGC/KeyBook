<svelte:options tag="device-new-svelte" />
<script lang="ts">
    import DeviceDetails from './DeviceDetails.svelte';

    let device = {
        "name": null,
        "identifier": null,
        "type": null,
    };
</script>
<main>
    <DeviceDetails
      bind:name={device.name}
      bind:identifier={device.identifier}
      bind:type={device.type}
      hidestatus=true />
    <div>
        <form action="/Device/Add" method="POST">
            {#each Object.keys(device) as key}
                <input type="hidden" name="{key}" id="{key}" value="{device[key]}" />
            {/each}
            <button>Add device</button>
        </form>
    </div>
</main>