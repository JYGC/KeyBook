export interface IDeviceListItemModel {
  id: string;
  deviceid: string;
  devicename: string;
  deviceidentifier: string;
  devicetype: string;
  personname: string;
}

export interface IEditDeviceModel {
  id: string;
  type: string;
  name: string;
  identifier: string;
  defunctreason: string;
  property: string;
}