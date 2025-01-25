import PocketBase, { type AuthModel } from "pocketbase";
import type { IBackendClient } from "$lib/interfaces";
import { localPocketbaseUrl } from "./api-constants";

export class BackendClient implements IBackendClient {
  private readonly __pb = new PocketBase();
  
  constructor() {
    this.__pb = new PocketBase(localPocketbaseUrl);
    this.__pb.authStore.loadFromCookie(document.cookie);
    this.__pb.authStore.onChange(() => {
      document.cookie = this.__pb.authStore.exportToCookie({ httpOnly: false });
    });
  }

  public get isTokenValid() {
    return this.__pb.authStore.isValid;
  }

  public logoutAsync = async () => {
    await this.__pb.collection('users').authRefresh();
    this.__pb.authStore.clear();
  };

  public authRefresh = async () => await this.__pb.collection('users').authRefresh();

  public get pb() {
    return this.__pb;
  };

  public get loggedInUser(): AuthModel {
    return this.__pb.authStore.model;
  };
}