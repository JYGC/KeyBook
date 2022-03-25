<svelte:options tag="person-new-svelte" />
<script lang="ts">
    import PersonDetails from './PersonDetails.svelte';

    let person = {
        "name": null,
        "type": "0"
    };
</script>
<main>
    <PersonDetails
      bind:name={person.name}
      bind:type={person.type}
      hideisgone=true/>
    <div>
        <form action="/Person/Add" method="POST">
            {#each Object.keys(person) as key}
                <input type="hidden" name="{key}" id="{key}" value="{person[key]}" />
            {/each}
            <button>Add person</button>
        </form>
    </div>
</main>