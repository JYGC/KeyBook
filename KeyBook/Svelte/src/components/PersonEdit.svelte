<svelte:options tag="person-edit-svelte" />
<script lang="ts">
	export let personjson;

    let person = JSON.parse(personjson);
    import PersonDetails from './PersonDetails.svelte';
</script>
<main>
    <button on:click="window.history.back()">Back</button>
    <PersonDetails bind:name={person.name} bind:type={person.type} bind:isgone={person.isGone} disabletype=true />
    {#if person.personDevices && person.personDevices.length > 0}
        <div>
            <label for="name">Devices Held:</label>
        </div>
        <div class="persns-devices">
            <table>
                <tr>
                    <th>Device Name</th>
                    <th>Device Identifier</th>
                    <th>Type</th>
                    <th></th>
                </tr>
                {#each person.personDevices as personDevice}
                    <tr>
                        <td>{personDevice.device.name}</td>
                        <td>{personDevice.device.identifier}</td>
                        <td>{personDevice.device.type}</td>
                        <td>
                            <button on:click={
                                () => window.location.href = "/Device/Edit/" + personDevice.device.id + "?fromPersonDetailsPersonId=" + person.id
                            }>Details</button>
                        </td>
                    </tr>
                {/each}
            </table>
        </div>
    {/if}
    <div>
        <form action="/Person/Save" method="POST">
            <input type="hidden" name="personid" id="personid" value={person.id} />
            <input type="hidden" name="name" id="name" value="{person.name}" />
            <input type="hidden" name="isgone" id="isgone" value="{person.isGone}" />
            <button>Save person</button>
        </form>
    </div>
</main>