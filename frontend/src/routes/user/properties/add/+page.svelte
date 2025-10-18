<script lang="ts">
	import { goto } from "$app/navigation";
	import { getBackendClient } from "$lib/api/backend-client";
	import PropertyEditor from "$lib/components/property/PropertyEditor.svelte";
	import type { IEditPropertyModel } from "$lib/models/property-models";
	import { Button } from "carbon-components-svelte";

  const gotoPropertyList = () => {
    goto("/user/properties/list");
  }

  const savePropertyActionAsync = async (changedProperty: IEditPropertyModel) => {
    try {
      const backendClient = getBackendClient();
      if (backendClient.authStore.record === null) {
        throw new Error("Cannot find loggedInUser.");
      }
      changedProperty.owners = [ backendClient.authStore.record.id ];
      await backendClient.collection("properties").create<IEditPropertyModel>(changedProperty);
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