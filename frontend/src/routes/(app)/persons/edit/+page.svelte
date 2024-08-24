<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import PersonEditor from "$lib/components/person/PersonEditor.svelte";
	import { getPersonContext } from "$lib/contexts/person-context.svelte";
	import type { IEditPersonDto } from "$lib/dtos/person-dtos";
	import { Button, Tile } from "carbon-components-svelte";

  const personContext = getPersonContext();

  const backendClient = new BackendClient();

  let personAsync = $derived.by<Promise<IEditPersonDto | null>>(async () => {
    try {
      return await backendClient.pb.collection("persons").getOne<IEditPersonDto>(
        personContext.selectedPersonId,
        { fields: "id,name,type,property" },
      );
    } catch (ex) {
      alert(ex);
      return null;
    }
  });

  const gotoPropertyPersonList = () => {
    goto("/persons/list/property");
  };

  const savePersonActionAsync = async (changedPerson: IEditPersonDto) => {
    try {
      await backendClient.pb.collection("persons").update(changedPerson.id, {
        type: changedPerson.type,
        name: changedPerson.name,
      });
      gotoPropertyPersonList();
    } catch (ex) {
      alert(ex);
    }
  };

  const deletePersonActionAsync = async (person: IEditPersonDto) => {
    try {
      await backendClient.pb.collection("persons").delete(person.id);
      gotoPropertyPersonList();
    } catch (ex) {
      alert(ex);
    }
  };
</script>

<Button onclick={gotoPropertyPersonList}>Back</Button>

{#await personAsync}
  <Tile>...getting person details</Tile>
{:then person}
  {#if person === null}
    {history.back()}
  {:else}
    <PersonEditor
      person={person}
      isAdd={false}
      savePersonAction={savePersonActionAsync}
      deletePersonAction={deletePersonActionAsync}
    />
  {/if}
{:catch error}
  {error}
{/await}