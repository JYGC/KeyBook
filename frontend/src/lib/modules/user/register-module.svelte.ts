import PocketBase from "pocketbase";
import type { IRegisterApi } from "$lib/interfaces";

export class RegisterModule implements IRegisterApi {
  private readonly __backendClient: PocketBase;
  
  public name = $state<string>("");
  public email = $state<string>("");
  public password = $state<string>("");
  public passwordConfirm = $state<string>("");

  constructor(backendClient: PocketBase) {
    this.__backendClient = backendClient;
  }

  private __error = $state<string>("");

  public callApi = async (): Promise<boolean> => {
    try {
      await this.__backendClient.collection("users").create({
        name: this.name,
        password: this.password,
        passwordConfirm: this.passwordConfirm,
        email: this.email,
        emailVisibility: false,
      });
      return true;
    } catch (ex) {
      this.__error = JSON.stringify(ex);
      return false;
    }
  };

  public get error() {
    return this.__error;
  }
}