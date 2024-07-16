import PocketBase from "pocketbase";
import type { AddPropertyDeviceAndHistoriesDTO } from "./dtos";

export interface IBackendClient {
  isTokenValid: boolean;
  logoutAsync: () => void;
  authRefresh: () => void;
  get pb(): PocketBase;
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
  callApi: () => Promise<boolean>;
  get error(): string;
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