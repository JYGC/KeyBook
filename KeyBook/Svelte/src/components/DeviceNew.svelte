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
    <button on:click="{() => window.history.back()}">Back</button>
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
            <input type="hidden" name="name" id="name" value="{device.name}" />
            <input type="hidden" name="identifier" id="identifier" value="{device.identifier}" />
            <input type="hidden" name="type" id="type" value="{device.type}" />
            <button>Add device</button>
        </form>
    </div>
</main>