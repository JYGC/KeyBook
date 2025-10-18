<script lang="ts">
	import { getBackendClient } from "$lib/api/backend-client";
	import PropertyEditor from "$lib/components/property/PropertyEditor.svelte";
	import ConfirmButtonAndDialog from "$lib/components/shared/ConfirmButtonAndDialog.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import { PropertyUpdateEditorModule } from "$lib/modules/property/property-update-editor-module.svelte";
	import { Button, Tile } from "carbon-components-svelte";

  const propertyContext = getPropertyContext();

  const backendClient = getBackendClient();

  const goBack = () => {
    history.back();
  };
  
  const propertyAddEditorModule = new PropertyUpdateEditorModule(
    backendClient,
    propertyContext,
    goBack,
  );
</script>

<Button onclick={goBack}>Back</Button>

<PropertyEditor propertyEditorModule={propertyAddEditorModule}>
  {#snippet deleteButton(deleteActionButtonClick: () => null)}
    <ConfirmButtonAndDialog
      submitAction={deleteActionButtonClick}
      buttonText="Delete"
      bodyMessage="Are you sure you want to delete this property?"
    />
  {/snippet}
</PropertyEditor>