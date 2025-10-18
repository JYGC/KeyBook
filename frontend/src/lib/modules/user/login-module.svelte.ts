import PocketBase from "pocketbase";
import type { ILoginModule } from "$lib/modules/interfaces";

export class LoginModule implements ILoginModule {
  private readonly __backendClient: PocketBase;
  
  public email = $state<string>("");
  public password = $state<string>("");

  constructor(backendClient: PocketBase) {
    this.__backendClient = backendClient;
  }

  public callApi = async (): Promise<string> => {
    await this.__backendClient.collection("users").authWithPassword(this.email, this.password);
    return this.__backendClient.authStore.exportToCookie({ httpOnly: false });
  };
}