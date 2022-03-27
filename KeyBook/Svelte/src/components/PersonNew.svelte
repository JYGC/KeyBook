<svelte:options tag="person-new-svelte" />
<script lang="ts">
    import PersonDetails from './PersonDetails.svelte';

    let person = {
        "name": null,
        "type": "0"
    };
</script>
<main>
    <button on:click="{() => window.history.back()}">Back</button>
    <PersonDetails
      bind:name={person.name}
      bind:type={person.type}
      hideisgone=true/>
    <div>
        <form action="/Person/Add" method="POST">
            <input type="hidden" name="name" id="name" value="{person.name}" />
            <input type="hidden" name="type" id="type" value="{person.type}" />
            <button>Add person</button>
        </form>
    </div>
</main>