import PocketBase, { type AuthModel } from "pocketbase";
import type { IBackendClient } from "$lib/interfaces";
import { PUBLIC_POCKETBASE_URL, PUBLIC_UNSECURE_COOKIE } from '$env/static/public';

export class BackendClient implements IBackendClient {
  private readonly __pb = new PocketBase();

  constructor() {
    this.__pb = new PocketBase(PUBLIC_POCKETBASE_URL);
    this.__pb.authStore.loadFromCookie(document.cookie);
    this.__pb.authStore.onChange(() => {
      document.cookie = this.__pb.authStore.exportToCookie({
        httpOnly: false,
        secure: PUBLIC_UNSECURE_COOKIE === undefined ||
          PUBLIC_UNSECURE_COOKIE === null ||
          PUBLIC_UNSECURE_COOKIE !== 'true'
      });
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

export const getBackendClient = () => {
  const pb = new PocketBase(PUBLIC_POCKETBASE_URL);
  pb.authStore.loadFromCookie(document.cookie);
  pb.authStore.onChange(() => {
    document.cookie = pb.authStore.exportToCookie({
      httpOnly: false,
      secure: PUBLIC_UNSECURE_COOKIE === undefined ||
        PUBLIC_UNSECURE_COOKIE === null ||
        PUBLIC_UNSECURE_COOKIE !== 'true'
    });
  });
  return pb;
}