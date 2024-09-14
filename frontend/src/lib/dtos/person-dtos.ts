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

export interface IPersonDeviceExpandPersonDevicePersonEditModel {
  id: string;
  expand: {
    person: IPersonIdNameTypeModel
  };
}