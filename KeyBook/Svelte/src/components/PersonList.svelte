<svelte:options tag="person-list-svelte" />
<script lang="ts">
    export let personsjson;
    let persons = JSON.parse(personsjson);

    let personTypes = {};
    import { GetPersonTypesAPI } from '../api/person';
    let getPersonTypesAPI = new GetPersonTypesAPI();
    getPersonTypesAPI.successCallback = r => { personTypes = r.data; };
    getPersonTypesAPI.failedCallback = e => { alert('error: ' + e); /* Add error handling later */ };
    getPersonTypesAPI.call();
</script>

<main>
    <div>
        <!-- FIX LATER: Back button goes to device if device was saved from person edit -->
        <button on:click={() => window.location.href = "/Person/New"}>Add person</button>
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
                    <td>{personTypes[person.type]}</td>
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
