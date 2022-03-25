<svelte:options tag="device-list-svelte" />
<script lang="ts">
	export let devicesjson;
    let devices = JSON.parse(devicesjson);
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
                    <td>{device.type}</td>
                    <td>
                        <button on:click={() => window.location.href = "/Device/Edit/" + device.id}>Details</button>
                    </td>
                </tr>
            {/each}
        </table>
    </div>
</main>
