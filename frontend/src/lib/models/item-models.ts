export interface IItemModel {
  id: string;
  name: string;
  description: string;
}

export interface IEntryDeviceModel {
  id: string;
  item: string;
  deviceType: string;
  identifier: string;
  defunctReason: string;
}

export interface IPropertyItemModel {
  id: string;
  item: string;
  property: string;
}

export interface IPersonItemModel {
  id: string;
  person: string;
  item: string;
}

export interface IItemListItemModel {
  id: string;
  name: string;
  description: string;
}

export interface IEditItemModel {
  id: string;
  name: string;
  description: string;
}

export interface IEditEntryDeviceModel {
  id: string;
  item: string;
  deviceType: string;
  identifier: string;
  defunctReason: string;
}
