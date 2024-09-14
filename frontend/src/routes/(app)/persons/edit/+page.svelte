<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import PersonEditor from "$lib/components/person/PersonEditor.svelte";
	import { getPersonContext } from "$lib/contexts/person-context.svelte";
	import { PersonUpdateEditorModule } from "$lib/modules/person/person-update-editor-module.svelte";
	import { Button, Tile } from "carbon-components-svelte";

  const personContext = getPersonContext();

  const backendClient = new BackendClient();

  const gotoPropertyPersonList = () => {
    goto("/persons/list/property");
  };

  const personEditorModule = new PersonUpdateEditorModule(
    backendClient,
    personContext,
    gotoPropertyPersonList
  );
</script>

<Button onclick={gotoPropertyPersonList}>Back</Button>
<PersonEditor personEditorModule={personEditorModule} />
