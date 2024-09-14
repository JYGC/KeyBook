export interface IDeviceListItemModel {
  id: string;
  type: string;
  name: string;
  identifier: string;
}

export interface IEditDeviceModel {
  id: string;
  type: string;
  name: string;
  identifier: string;
  defunctreason: string;
  property: string;
}