<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import DeviceEdit from "$lib/components/device/DeviceEditor.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import { type IEditDeviceDto } from "$lib/dtos/device-dtos";
	import { Button } from "carbon-components-svelte";

  const propertyContext = getPropertyContext();

const gotoPropertyList = () => {
  goto("/devices/list/property");
}

  const saveDeviceActionAsync = async (changedDevice: IEditDeviceDto) => {
    try {
      const backendClient = new BackendClient();
      changedDevice.property = propertyContext.selectedPropertyId;
      await backendClient.pb.collection("devices").create<IEditDeviceDto>(changedDevice);
      gotoPropertyList();
    } catch (ex) {
      alert(ex);
    }
  };
</script>
<Button onclick={gotoPropertyList}>Back</Button>
<DeviceEdit
  device={{} as IEditDeviceDto}
  isAdd={true}
  saveDeviceAction={saveDeviceActionAsync}
/>