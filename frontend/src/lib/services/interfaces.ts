import type { IDeviceListItemModel } from "$lib/models/person-device-models";
import type { AuthModel } from "pocketbase";

export interface IAuthService {
  get isTokenValid(): boolean;
  logoutAsync: () => void;
  authRefresh: () => void;
  get loggedInUser(): AuthModel;
}

export interface IDeviceListViewService {
  getForDeviceListAsync: (selectedPropertyId: string) => Promise<IDeviceListItemModel[]>;
}