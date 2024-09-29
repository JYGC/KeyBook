import PocketBase, { type AuthModel } from "pocketbase";
import type { AddPropertyDeviceAndHistoriesModel } from "./models/data-import-models";
import type { IDeviceListItemModel, IEditDeviceModel } from "./models/device-models";
import type { IEditPersonModel, IPersonDeviceExpandPersonDevicePersonEditModel, IPersonIdNameTypeModel } from "./models/person-models";

export interface IBackendClient {
  isTokenValid: boolean;
  logoutAsync: () => void;
  authRefresh: () => void;
  get pb(): PocketBase;
  get loggedInUser(): AuthModel;
}

export interface ILoginApi {
  email: string;
  password: string;
  callApi: () => Promise<string>;
}

export interface IRegisterApi {
  name: string;
  email: string;
  password: string;
  passwordConfirm: string;
  callApi: () => Promise<boolean>;
  get error(): string;
}

export interface IDeviceEditorModule {
  deviceAsync: Promise<IEditDeviceModel | null>;
  get isAdd(): boolean;
  deviceStatusTextAsync: Promise<string>;
  getSaveDeviceAction: () => ((device: IEditDeviceModel) => void);
  getDeleteDeviceAction: () => ((device: IEditDeviceModel) => void) | null;
  callBackAction: () => void;
}

export interface IPersonEditorModule {
  personAsync: Promise<IEditPersonModel | null>;
  get isAdd(): boolean;
  getSavePersonAction: () => ((device: IEditPersonModel) => void);
  getDeletePersonAction: () => ((device: IEditPersonModel) => void) | null;
  callBackAction: () => void;
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

export interface IUploadCsvApi {
  callApi: () => void;
}

export interface ICsvFileToObjectConverter {
  input: FileList | null;
  outputAsync: Promise<AddPropertyDeviceAndHistoriesModel | null>;
}

export interface IDeviceListDataService {
  add: () => void;
}