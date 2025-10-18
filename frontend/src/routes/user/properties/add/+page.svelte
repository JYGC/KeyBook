<script lang="ts">
	import { PropertyAddEditorModule } from "$lib/modules/property/property-add-editor-module.svelte";
	import { getBackendClient } from "$lib/api/backend-client";
	import PropertyEditor from "$lib/components/property/PropertyEditor.svelte";
	import { Button } from "carbon-components-svelte";

  const goBack = () => {
    history.back();
  };

  const backendClient = getBackendClient();
  const propertyAddEditorModule = new PropertyAddEditorModule(
    backendClient,
    goBack,
  );
</script>

<Button onclick={goBack}>Back</Button>
<PropertyEditor
  propertyEditorModule={propertyAddEditorModule}
/>