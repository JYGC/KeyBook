<svelte:options tag="person-list-svelte" />
<script lang="ts">
    export let personsjson;
    let persons = JSON.parse(personsjson);
</script>

<main>
    <div>
        <button on:click={() => window.location.href = "/Person/New"}>Add device</button>
    </div>
    <div>
        <table>
            <tr>
                <th>Person Name</th>
                <th>IsGone</th>
                <th>Person Type</th>
                <th>Current Devices</th>
                <th></th>
            </tr>
            {#each persons as person}
                <tr>
                    <td>{person.name}</td>
                    <td>{person.isGone}</td>
                    <td>{person.type}</td>
                    <td>
                        {#each person.personDevices as personDevice}
                            {personDevice.device.name}<br />
                        {/each}
                    </td>
                    <td>
                        <button on:click={() => window.location.href = "/Person/Edit/" + person.id}>Details</button>
                    </td>
                </tr>
            {/each}
        </table>
    </div>
</main>
