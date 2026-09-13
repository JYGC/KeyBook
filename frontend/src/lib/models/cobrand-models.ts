export interface ICobrandModel {
  id: string;
  name: string;
}

export interface ICobrandAdminModel {
  id: string;
  user: string;
  cobrand: string;
  approved: boolean;
}

export interface ICobrandPropertyManagerModel {
  id: string;
  cobrand: string;
  property: string;
}

export interface ICobrandPropertyOwnerModel {
  id: string;
  cobrand: string;
  propertyOwner: string;
}
