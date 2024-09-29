import type { IDeviceListItemModel } from "./device-models";

export interface IPersonDeviceModel {
  person: string;
  device: string;
}

export interface IDeviceHeldExpandDeviceModel {
  id: string;
  expand: {
    device: IDeviceListItemModel
  }
}