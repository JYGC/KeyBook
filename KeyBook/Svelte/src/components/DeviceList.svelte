<svelte:options tag="device-list-svelte" />
<script lang="ts">
	export let devicesjson;
    let devices = JSON.parse(devicesjson);

    let deviceTypes = {};
    import { GetDeviceTypesAPI } from '../api/device';
    let getDeviceTypesAPI = new GetDeviceTypesAPI();
    getDeviceTypesAPI.successCallback = r => { deviceTypes = r.data; };
    getDeviceTypesAPI.failedCallback = e => { alert('error: ' + e); /* Add error handling later */ }
    getDeviceTypesAPI.call();
</script>

<main>
    <div>
        <button on:click={() => window.location.href = "/Device/New"}>Add device</button>
    </div>
    <div>
        <table>
            <tr>
                <th>Device Name</th>
                <th>Device Identifier</th>
                <!--<th>Status</th>-->
                <th>Type</th>
                <th></th>
            </tr>
            {#each devices as device}
                <tr>
                    <td>{device.name}</td>
                    <td>{device.identifier}</td>
                    <!--<td>{device.Status}</td>-->
                    <td>{deviceTypes[device.type]}</td>
                    <td>
                        <button on:click={() => window.location.href = "/Device/Edit/" + device.id}>Details</button>
                    </td>
                </tr>
            {/each}
        </table>
    </div>
</main>
