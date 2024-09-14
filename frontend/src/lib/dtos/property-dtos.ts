export interface IPropertyListItemModel {
  id: string;
  address: string;
}

export interface IEditPropertyModel {
  id: string;
  address: string;
  owners: Array<string>;
  managers: Array<string>;
}