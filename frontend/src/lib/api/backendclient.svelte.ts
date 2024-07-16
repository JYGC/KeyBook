import PocketBase from "pocketbase";
import { PUBLIC_POCKETBASE_URL } from "$env/static/public";
import type { IBackendClient } from "$lib/interfaces";

export class BackendClient implements IBackendClient {
  private readonly __pb = $state<PocketBase>(new PocketBase());
  
  constructor() {
    this.__pb = new PocketBase(PUBLIC_POCKETBASE_URL);
    this.__pb.authStore.loadFromCookie(document.cookie);
  }

  public isTokenValid = $derived(this.__pb.authStore.isValid);

  public logoutAsync = async () => {
    await this.__pb.collection('users').authRefresh();
    this.__pb.authStore.clear();
    // set past expiry to all cookies
    const cookies = document.cookie.split(';');
    for (let i = 0; i < cookies.length; i++) {
      document.cookie = cookies[i] + "=; expires="+ new Date(0).toUTCString();
    }
  };

  public authRefresh = async () => await this.__pb.collection('users').authRefresh();

  public get pb() {
    return this.__pb;
  };
}