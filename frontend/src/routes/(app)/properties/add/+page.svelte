<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import PropertyEditor from "$lib/components/property/PropertyEditor.svelte";
	import type { IEditPropertyDto } from "$lib/dtos/property-dtos";
	import { Button } from "carbon-components-svelte";

  const gotoPropertyList = () => {
    goto("/properties/list");
  }

  const savePropertyActionAsync = async (changedProperty: IEditPropertyDto) => {
    try {
      const backendClient = new BackendClient();
      if (backendClient.loggedInUser === null) {
        throw new Error("Cannot find loggedInUser.");
      }
      changedProperty.owners = backendClient.loggedInUser.id;
      await backendClient.pb.collection("properties").create<IEditPropertyDto>(changedProperty);
      gotoPropertyList();
    } catch (ex) {
      alert(ex);
    }
  };
</script>

<Button onclick={gotoPropertyList}>Back</Button>
<PropertyEditor
  property={{} as IEditPropertyDto}
  isAdd={true}
  savePropertyAction={savePropertyActionAsync}
/>