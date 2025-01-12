<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import PersonList from "$lib/components/person/PersonList.svelte";
	import { getPersonContext } from "$lib/contexts/person-context.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import { PersonListModule } from "$lib/modules/person/person-list-module.svelte";
	import { Button } from "carbon-components-svelte";

	const propertyContext = getPropertyContext();
  if (
    propertyContext.selectedPropertyId === null ||
    propertyContext.selectedPropertyId.trim() === ""
  ) {
    goto("/");
  }

	const personContext = getPersonContext();

  const backendClient = new BackendClient();
	const personListModule = new PersonListModule(backendClient, propertyContext);

const gotoPropertyList = () => {
	goto("/properties/list");
};

const gotoAddNewProperty = () => {
	goto("/persons/add");
};
</script>

<Button onclick={gotoPropertyList}>Back</Button>
<Button onclick={gotoAddNewProperty}>Add New Person</Button>

<PersonList
	personListModule={personListModule}
	propertyId={propertyContext.selectedPropertyId}
	bind:selectedPersonId={personContext.selectedPersonId}
/>