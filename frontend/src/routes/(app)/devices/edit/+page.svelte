<script lang="ts">
	import { BackendClient } from "$lib/api/backend-client";
	import DeviceEditor from "$lib/components/device/DeviceEditor.svelte";
	import DeviceHolderEditor from "$lib/components/persondevice/DeviceHolderEditor.svelte";
	import { getDeviceContext } from "$lib/contexts/device-context.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import { DeviceUpdateEditorModule } from "$lib/modules/device/device-update-editor-module.svelte";
	import { DeviceHolderEditorModule } from "$lib/modules/persondevice/device-holder-editor-module.svelte";
	import { Button } from 'carbon-components-svelte';

  const propertyContext = getPropertyContext();
  const deviceContext = getDeviceContext();

  const backendClient = new BackendClient();

  const goBack = () => {
    history.back();
  };

  const deviceUpdateEditorModule = new DeviceUpdateEditorModule(
    backendClient,
    deviceContext,
    goBack,
  );

  const deviceHolderEditorModule = new DeviceHolderEditorModule(
    backendClient,
    deviceUpdateEditorModule,
    propertyContext,
  );
</script>

<Button onclick={goBack}>Back</Button>

<DeviceEditor deviceEditorModule={deviceUpdateEditorModule}>
  <DeviceHolderEditor deviceHolderEditorModule={deviceHolderEditorModule} />
</DeviceEditor>
