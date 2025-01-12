<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import PersonEditor from "$lib/components/person/PersonEditor.svelte";
	import DeviceHoldingList from "$lib/components/persondevice/DeviceHoldingList.svelte";
	import { getDeviceContext } from "$lib/contexts/device-context.svelte";
	import { getPersonContext } from "$lib/contexts/person-context.svelte";
	import { PersonUpdateEditorModule } from "$lib/modules/person/person-update-editor-module.svelte";
	import { DeviceHoldingListModule } from "$lib/modules/persondevice/device-holding-list-module.svelte";
	import { Button } from "carbon-components-svelte";

  const personContext = getPersonContext();
  const deviceContext = getDeviceContext();

  const backendClient = new BackendClient();

  const gotoPropertyPersonList = () => {
    goto("/persons/list/property");
  };

  const personEditorModule = new PersonUpdateEditorModule(
    backendClient,
    personContext,
    gotoPropertyPersonList
  );

  const deviceHoldingListModule = new DeviceHoldingListModule(backendClient, personContext);
</script>

<Button onclick={gotoPropertyPersonList}>Back</Button>
<PersonEditor personEditorModule={personEditorModule}>
  <DeviceHoldingList
    deviceHoldingListModule={deviceHoldingListModule}
    bind:selectedDeviceId={deviceContext.selectedDeviceId}
  />
</PersonEditor>
