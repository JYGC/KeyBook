<script lang="ts">
	import { type IEditPersonDto } from '$lib/dtos/person-dtos';
	import { goto } from "$app/navigation";
	import PersonEditor from "$lib/components/person/PersonEditor.svelte";
	import { Button } from "carbon-components-svelte";
	import { BackendClient } from '$lib/api/backend-client';
	import { getPropertyContext } from '$lib/contexts/property-context.svelte';

  const propertyContext = getPropertyContext();

  const gotoPropertyPersonList = () => {
    goto("/persons/list/property");
  };

  const saveDeviceActionAsync = async (changedPerson: IEditPersonDto) => {
    try {
      const backendClient = new BackendClient();
      changedPerson.property = propertyContext.selectedPropertyId;
      await backendClient.pb.collection("persons").create<IEditPersonDto>(changedPerson);
      gotoPropertyPersonList();
    } catch (ex) {
      alert(ex);
    }
  };
</script>

<Button onclick={gotoPropertyPersonList}>Back</Button>
<PersonEditor
  person={{} as IEditPersonDto}
  isAdd={true}
  savePersonAction={saveDeviceActionAsync}
/>