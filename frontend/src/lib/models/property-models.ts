export interface IPropertyModel {
  id: string;
  address: string;
}

export interface IPropertyOwnerModel {
  id: string;
  property: string;
}

export interface IPersonPropertyOwnerModel {
  id: string;
  person: string;
  propertyOwner: string;
}

export interface ITenantModel {
  id: string;
  person: string;
  property: string;
}

export interface IHouseholdModel {
  id: string;
  person: string;
  property: string;
}

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