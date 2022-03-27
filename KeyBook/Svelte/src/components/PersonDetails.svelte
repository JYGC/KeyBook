<svelte:options tag="person-details-svelte" />
<script lang="ts">
    export let name;
    export let type;
    export let isgone;
    export let hideisgone = false;
    export let disabletype = false;

    import Axios from 'axios';

    let persontypes = {};
    
    Axios.get(`/Person/GetPersonTypes`).then(response => {
        persontypes = response.data;
        type = String(type) // select value cannot recognise numbers
    }).catch(e => {
        alert('error: ' + e); // Add error handling later
    });
</script>
<main>
    <div>
        <div>
            <label for="name">Person Name:</label>
        </div>
        <div>
            <input type="text" name="name" id="name" bind:value="{name}" />
        </div>
        {#if !hideisgone}
            <div>
                <label for="isgone">Person Is Gone:</label>
            </div>
            <div>
                <input type="checkbox" name="isgone" id="name" bind:value="{isgone}" />
            </div>
        {/if}
        <div>
            <label for="type">Person Type:</label>
        </div>
        <div>
            <select name="type" id="type" bind:value={type} disabled={disabletype}>
                {#each Object.keys(persontypes) as key}
                    <option value={key}>{persontypes[key]}</option>
                {/each}
            </select>
        </div>
    </div>
</main>