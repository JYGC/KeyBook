export interface IEditPersonModel {
  id: string;
  type: string;
  name: string;
  property: string;
}

export interface IPersonIdNameTypeModel {
  id: string;
  type: string;
  name: string;
}

export interface IPersonListItemModel {
  id: string;
  personid: string;
  personname: string;
  persontype: string;
  devicejsons: string;
}

export interface IPersonDeviceExpandPersonDevicePersonEditModel {
  id: string;
  expand: {
    person: IPersonIdNameTypeModel
  };
}