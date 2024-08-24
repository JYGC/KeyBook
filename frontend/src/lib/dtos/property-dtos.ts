export interface IPropertyListItemDto {
  id: string;
  address: string;
}

export interface IEditPropertyDto {
  id: string;
  address: string;
  owners: Array<string>;
  managers: Array<string>;
}