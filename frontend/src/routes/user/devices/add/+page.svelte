<script lang="ts">
	import { goto } from "$app/navigation";
	import { getBackendClient } from "$lib/api/backend-client";
	import DeviceEditor from "$lib/components/device/DeviceEditor.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import { DeviceAddEditorModule } from "$lib/modules/device/device-add-editor-module.svelte";
	import { Button } from "carbon-components-svelte";

  const propertyContext = getPropertyContext();

  const backendClient = getBackendClient();

  const gotoPropertyList = () => {
    goto("/user/devices/list/property");
  }

  const deviceAddEditorModule = new DeviceAddEditorModule(
    backendClient,
    propertyContext,
    gotoPropertyList,
  );
</script>

<Button onclick={gotoPropertyList}>Back</Button>
<DeviceEditor deviceEditorModule={deviceAddEditorModule} />