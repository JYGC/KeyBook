<script lang="ts">
	import { goto } from "$app/navigation";
	import PersonEditor from "$lib/components/person/PersonEditor.svelte";
	import { Button } from "carbon-components-svelte";
	import { BackendClient } from '$lib/api/backend-client';
	import { getPropertyContext } from '$lib/contexts/property-context.svelte';
	import { PersonAddEditorModule } from '$lib/modules/person/person-add-editor-module.svelte';

  const propertyContext = getPropertyContext();

  const backendClient = new BackendClient();

  const gotoPropertyPersonList = () => {
    goto("/user/persons/list/property");
  };

  const personEditorModule = new PersonAddEditorModule(
    backendClient,
    propertyContext,
    gotoPropertyPersonList
  );
</script>

<Button onclick={gotoPropertyPersonList}>Back</Button>
<PersonEditor personEditorModule={personEditorModule} />