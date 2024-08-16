<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import PersonList from "$lib/components/person/PersonList.svelte";
	import { getPersonContext } from "$lib/contexts/person-context.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import type { IPersonListItemDto } from "$lib/dtos/person-dtos";
	import { Button, Tile } from "carbon-components-svelte";

	const propertyContext = getPropertyContext();
  if (
    propertyContext.selectedPropertyId === null ||
    propertyContext.selectedPropertyId.trim() === ""
  ) {
    goto("/");
  }

	const personContext = getPersonContext();

  const backendClient = new BackendClient();

	let personListAsync = $derived.by<Promise<IPersonListItemDto[]>>(async () => {
		try {
			const response = await backendClient.pb.collection("persons").getList<IPersonListItemDto>(
				1,
				50,
				{
					filter: `property.id = "${propertyContext.selectedPropertyId}"`,
					fields: "id,type,name",
				}
			);
			return response.items;
		} catch (ex) {
			alert(ex);
			return [];
		}
	});

const gotoPropertyList = () => {
	goto("/properties/list");
};
</script>

<Button onclick={gotoPropertyList}>Back</Button>

{#await personListAsync}
<Tile>...getting persons</Tile>
{:then personList}
	<PersonList
		personList={personList}
		propertyId={propertyContext.selectedPropertyId}
		bind:selectedPersonId={personContext.selectedPersonId}
	/>
{:catch error}
  {error}
{/await}