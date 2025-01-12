import { getContext, setContext } from "svelte";

export class DeviceContext {
  public selectedDeviceId = $state<string>("");
}

export const setDeviceContext = () => {
  const deviceContext = new DeviceContext();
  setContext<DeviceContext>("deviceContext", deviceContext)
};

export const getDeviceContext = () => getContext<DeviceContext>("deviceContext");