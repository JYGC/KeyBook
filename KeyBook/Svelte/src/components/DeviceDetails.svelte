<svelte:options tag="device-details-svelte" />
<script lang="ts">
    export let name;
    export let identifier;
    export let status;
    export let type;
    export let hidestatus = false;
    export let disabletype = false;

    let deviceTypes = {};

    import { getDeviceTypes } from '../api/get-types';
    getDeviceTypes(r => {
        deviceTypes = r.data;
        type = String(type) // select value cannot recognise numbers
    }, e => { alert('error: ' + e); /* Add error handling later */ });
</script>
<main>
    <div>
        <div>
            <label for="name">Device Name:</label>
        </div>
        <div>
            <input type="text" name="name" id="name" bind:value={name} />
        </div>
        <div>
            <label for="identifier">Device Identifier:</label>
        </div>
        <div>
            <input type="text" name="identifier" id="identifier" bind:value={identifier} />
        </div>
        {#if !hidestatus}
            <div>
                <label for="status">Device Status:</label>
            </div>
            <div>
                <input type="text" name="status" id="status" bind:value={status} />
            </div>
        {/if}
        <div>
            <label for="type">Device Type:</label>
        </div>
        <div>
            <select name="type" id="type" bind:value={type} disabled={disabletype}>
                {#each Object.keys(deviceTypes) as key}
                    <option value={key}>{deviceTypes[key]}</option>
                {/each}
            </select>
        </div>
    </div>
</main>
