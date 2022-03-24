<svelte:options tag="device-edit-svelte" />
<script lang="ts">
    export let devicejson;

    import DeviceDetails from './DeviceDetails.svelte';
    let device = JSON.parse(devicejson);
</script>
<main>
    <DeviceDetails
      bind:name={device.name}
      bind:identifier={device.identifier}
      bind:status={device.status}
      bind:type={device.type}
      disabletype=true />
    <div>
        <label for="person">Person holding device:</label>
    </div>
    <div>
        <select name="cars" id="cars">
            <option value="volvo">Volvo</option>
            <option value="saab">Saab</option>
            <option value="mercedes">Mercedes</option>
            <option value="audi">Audi</option>
        </select>
    </div>
</main>