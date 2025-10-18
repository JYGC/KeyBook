import type { IEditDeviceModel } from "$lib/models/device-models";
import type { IDeviceHistoryListItem, IDeviceListItemModel, IPersonListItemModel } from "$lib/models/person-device-models";
import type { IEditPersonModel, IPersonDeviceExpandPersonDevicePersonEditModel, IPersonIdNameTypeModel } from "$lib/models/person-models";
import type { IEditPropertyModel, IPropertyListItemModel } from "$lib/models/property-models";

export interface IPropertyListModule {
  propertyListAsync: Promise<IPropertyListItemModel[]>;
}

export interface IPropertyEditorModule {
  propertyAsync: Promise<IEditPropertyModel | null>;
  get isAdd(): boolean;
  savePropertyAction: () => ((changedProperty: IEditPropertyModel) => void);
  getDeletePropertyAction: () => ((property: IEditPropertyModel) => void) | null;
  callBackAction: () => void;
}

export interface IDeviceEditorModule {
  deviceAsync: Promise<IEditDeviceModel | null>;
  get isAdd(): boolean;
  deviceStatusTextAsync: Promise<string>;
  getSaveDeviceAction: () => ((device: IEditDeviceModel) => void);
  getDeleteDeviceAction: () => ((device: IEditDeviceModel) => void) | null;
  callBackAction: () => void;
}

export interface IDeviceListModule {
  deviceListAsync: Promise<IDeviceListItemModel[]>;
}

export interface IPersonEditorModule {
  personAsync: Promise<IEditPersonModel | null>;
  get isAdd(): boolean;
  getSavePersonAction: () => ((device: IEditPersonModel) => void);
  getDeletePersonAction: () => ((device: IEditPersonModel) => void) | null;
  callBackAction: () => void;
}

export interface IPersonListModule {
  personListAsync: Promise<IPersonListItemModel[]>;
}

export interface IDeviceHolderEditorModule {
  personDeviceExpandPersonDevicePersonAsync:
    Promise<IPersonDeviceExpandPersonDevicePersonEditModel | null>;
  availablePersonsAsync: Promise<IPersonIdNameTypeModel[]>;
  replaceDeviceHolderActionAsync: () => void;
  currentDeviceHolderNameAsync: Promise<string>;
  selectedDeviceHolderId: string;
}

export interface IDeviceHoldingListModule {
  deviceHoldingListOfPersonAsync: Promise<IDeviceListItemModel[]>;
}

export interface IDeviceHistoryListModule {
  deviceHistoryListAsync: Promise<IDeviceHistoryListItem[] | null>;
}

export interface IDeviceHistoryListUpdaterModule {
  updateTriggerState: unknown;
}

export interface ILoginModule {
  email: string;
  password: string;
  callApi: () => Promise<string>;
}

export interface IRegisterModule {
  name: string;
  email: string;
  password: string;
  passwordConfirm: string;
  callApi: () => Promise<boolean>;
  get error(): string;
}