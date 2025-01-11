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
  deviceid: string;
  devicename: string;
}

export interface IPersonListItemModel {
  id: string;
  personid: string;
  personname: string;
  persontype: string;
  holdingdevicejsons: IHoldingDeviceIdName[];
}
export interface IDeviceListItemModel {
  id: string;
  deviceid: string;
  devicename: string;
  deviceidentifier: string;
  devicetype: string;
  personname: string;
}

