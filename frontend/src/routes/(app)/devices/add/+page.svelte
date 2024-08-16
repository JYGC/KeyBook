<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import DeviceEdit from "$lib/components/device/DeviceEditor.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import { type IEditDeviceDto } from "$lib/dtos/device-dtos";
	import { Button } from "carbon-components-svelte";

  const propertyContext = getPropertyContext();

  const saveDeviceActionAsync = async (changedDevice: IEditDeviceDto) => {
    try {
      const backendClient = new BackendClient();
      await backendClient.pb.collection("devices").create<IEditDeviceDto>(changedDevice);
    } catch (ex) {
      alert(ex);
    }
  };

	const gotoPropertyListAsync = async () => {
    goto("/devices/listinproperty");
	}
</script>
<Button onclick={gotoPropertyListAsync}>Back</Button>
<DeviceEdit
  device={{ property: propertyContext.selectedPropertyId } as IEditDeviceDto}
  isAdd={true}
  saveDeviceAction={saveDeviceActionAsync}
/>