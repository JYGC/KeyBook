<script lang="ts">
	import { BackendClient } from "$lib/api/backend-client";
	import DeviceEditor from "$lib/components/device/DeviceEditor.svelte";
	import DeviceHistoryList from "$lib/components/persondevice/DeviceHistoryList.svelte";
	import DeviceHolderEditor from "$lib/components/persondevice/DeviceHolderEditor.svelte";
	import ConfirmButtonAndDialog from "$lib/components/shared/ConfirmButtonAndDialog.svelte";
	import { getDeviceContext } from "$lib/contexts/device-context.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import { DeviceUpdateEditorModule } from "$lib/modules/device/device-update-editor-module.svelte";
	import { DeviceHistoryListModule } from "$lib/modules/persondevice/device-history-list-module.svelte";
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

  const deviceHistoryListModule = new DeviceHistoryListModule(
    backendClient,
    deviceContext,
    deviceHolderEditorModule,
  )
</script>

<Button onclick={goBack}>Back</Button>

<DeviceEditor deviceEditorModule={deviceUpdateEditorModule}>
  <DeviceHolderEditor deviceHolderEditorModule={deviceHolderEditorModule} />
  {#snippet deleteButton(deleteActionButtonClick: () => null)}
    <ConfirmButtonAndDialog
      submitAction={deleteActionButtonClick}
      buttonText="Delete"
      bodyMessage="Are you sure you want to delete this device?"
    />
  {/snippet}
</DeviceEditor>
<br />
<br />
<br />
<br />
<DeviceHistoryList deviceHistoryListModule={deviceHistoryListModule} />