<svelte:options tag="device-edit-svelte" />
<script lang="ts">
    export let devicepersonlistjson;

    let devicepersonlist = JSON.parse(devicepersonlistjson);
    import DeviceDetails from './DeviceDetails.svelte';

    let personId;
    if (devicepersonlist.device.personDevice != null) personId = devicepersonlist.device.personDevice.personId;
    let saveableDetails = ["name", "identifier", "status", "type"]
</script>
<main>
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
            {#each devicepersonlist.personList as person}
                <option value="{person.id}">{person.name}</option>
            {/each}
        </select>
    </div>
    <div>
        <form action="/Device/Save" method="POST">
            <input type="hidden" name="deviceid" id="deviceid" value={devicepersonlist.device.id} />
            {#each saveableDetails as saveableDetail}
                <input type="hidden" name="{saveableDetail}" id="{saveableDetail}" value="{devicepersonlist.device[saveableDetail]}" />
            {/each}
            <input type="hidden" name="personid" id="personid" value={personId} />
            <button>Save device</button>
        </form>
    </div>
    {devicepersonlist.fromPersonDetailsPersonId}
</main>