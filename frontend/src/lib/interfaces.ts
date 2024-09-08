import PocketBase, { type AuthModel } from "pocketbase";
import type { AddPropertyDeviceAndHistoriesDTO } from "./dtos/data-import-dtos";
import type { IEditDeviceDto } from "./dtos/device-dtos";
import type { IPersonDeviceExpandPersonDevicePersonEditModel, IPersonIdNameTypeModel } from "./dtos/person-dtos";

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
  deviceAsync: Promise<IEditDeviceDto | null>;
  get isAdd(): boolean;
  deviceStatusTextAsync: Promise<string>;
  getSaveDeviceAction: () => ((device: IEditDeviceDto) => void);
  getDeleteDeviceAction: () => ((device: IEditDeviceDto) => void) | null;
  callBackAction: () => void;
}

export interface IDeviceHolderEditorModule {
  personDeviceExpandPersonDevicePersonAsync:
    Promise<IPersonDeviceExpandPersonDevicePersonEditModel | null>;
  availablePersonsAsync: Promise<IPersonIdNameTypeModel[]>;
  replaceDeviceHolderActionAsync: (selectedDeviceHolderId: string) => void;
}

export interface IUploadCsvApi {
  callApi: () => void;
}

export interface ICsvFileToObjectConverter {
  input: FileList | null;
  outputAsync: Promise<AddPropertyDeviceAndHistoriesDTO | null>;
}

export interface IDeviceListDataService {
  add: () => void;
}