export interface IPersonDeviceModel {
  person: string;
  device: string;
  property: string;
}

export interface IDeviceHeldExpandDeviceModel {
  id: string;
  expand: {
    device: IDeviceListItemModel
  }
}

export interface IHoldingDeviceIdName {
  deviceId: string;
  deviceName: string;
}

export interface IPersonListItemModel {
  id: string;
  personId: string;
  personName: string;
  personType: string;
  holdingDeviceJsons: IHoldingDeviceIdName[];
}
export interface IDeviceListItemModel {
  id: string;
  deviceId: string;
  deviceName: string;
  deviceIdentifier: string;
  deviceType: string;
  personName: string;
}

export interface IDeviceHistoryListItem {
  id: string;
  deviceId: string;
  description: string;
  created: string;
}