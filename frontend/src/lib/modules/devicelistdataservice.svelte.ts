import type { IDeviceListDataService } from "$lib/interfaces";

export class DeviceListDataService implements IDeviceListDataService {
  devices = $state<string[]>([]);
  private mustReload = $state<boolean>(true);

  constructor() {
    $effect(() => {
      if (this.mustReload) {
        this.devices.push(`${this.mustReload}`);
        this.mustReload = false;
      }
    });
  }

  add = () => {
    this.mustReload = true;
  }
};