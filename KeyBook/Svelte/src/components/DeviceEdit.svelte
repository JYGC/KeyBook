<svelte:options tag="device-edit-svelte" />
<script lang="ts">
    export let devicepersonlistjson;

    let devicepersonlist = JSON.parse(devicepersonlistjson);
    import DeviceDetails from './DeviceDetails.svelte';
    
    let personNames = {}
    import { GetPersonNamesAPI } from '../api/person';
    let getPersonNamesAPI = new GetPersonNamesAPI();
    getPersonNamesAPI.successCallback = r => { personNames = r.data; };
    getPersonNamesAPI.failedCallback = e => { alert('error: ' + e); /* Add error handling later */ }
    getPersonNamesAPI.call();

    let personId;
    if (devicepersonlist.device.personDevice != null) personId = devicepersonlist.device.personDevice.personId;
</script>
<main>
    <button on:click="{() => window.history.back()}">Back</button>
    <DeviceDetails
      bind:name={devicepersonlist.device.name}
      bind:identifier={devicepersonlist.device.identifier}
      bind:status={devicepersonlist.device.status}
      bind:type={devicepersonlist.device.type}
      disabletype=true />
    <div>
        <label for="person">Person holding device:</label>
    </div>
    <div>
        <select name="person" id="person" bind:value={personId}>
            <option value="">Not Used</option>
            {#each Object.keys(personNames) as personId}
                <option value="{personId}">{personNames[personId]}</option>
            {/each}
        </select>
    </div>
    <div>
        <form action="/Device/Save" method="POST">
            <input type="hidden" name="deviceid" id="deviceid" value={devicepersonlist.device.id} />
            <input type="hidden" name="name" id="name" value={devicepersonlist.device.name} />
            <input type="hidden" name="identifier" id="identifier" value={devicepersonlist.device.identifier} />
            <input type="hidden" name="status" id="status" value={devicepersonlist.device.status} />
            <input type="hidden" name="type" id="type" value={devicepersonlist.device.type} />
            <input type="hidden" name="personid" id="personid" value={personId} />
            <input type="hidden" name="frompersondetailspersonid" id="frompersondetailspersonid" value={devicepersonlist.fromPersonDetailsPersonId} />
            <button>Save device</button>
        </form>
    </div>
</main>