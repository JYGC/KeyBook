import type { IAuthService } from "$lib/services/interfaces";
import PocketBase, { type AuthRecord } from "pocketbase";

export class AuthService implements IAuthService {
  private readonly __pocketbase: PocketBase;
  
  constructor(pocketbase: PocketBase) {
    this.__pocketbase = pocketbase;
  }

  get isTokenValid(): boolean {
    return this.__pocketbase.authStore.isValid;
  }

  public logoutAsync = async () => {
    this.__pocketbase.authStore.clear();
  };

  public authRefresh = async () => await this.__pocketbase.collection('users').authRefresh();

  get loggedInUser(): AuthRecord {
    return this.__pocketbase.authStore.record;
  }
}