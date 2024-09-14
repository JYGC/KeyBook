<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import PropertyEditor from "$lib/components/property/PropertyEditor.svelte";
	import type { IEditPropertyModel } from "$lib/dtos/property-dtos";
	import { Button } from "carbon-components-svelte";

  const gotoPropertyList = () => {
    goto("/properties/list");
  }

  const savePropertyActionAsync = async (changedProperty: IEditPropertyModel) => {
    try {
      const backendClient = new BackendClient();
      if (backendClient.loggedInUser === null) {
        throw new Error("Cannot find loggedInUser.");
      }
      changedProperty.owners = backendClient.loggedInUser.id;
      await backendClient.pb.collection("properties").create<IEditPropertyModel>(changedProperty);
      gotoPropertyList();
    } catch (ex) {
      alert(ex);
    }
  };
</script>

<Button onclick={gotoPropertyList}>Back</Button>
<PropertyEditor
  property={{} as IEditPropertyModel}
  isAdd={true}
  savePropertyAction={savePropertyActionAsync}
/>