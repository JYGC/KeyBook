import type { IBackendClient, ILoginApi } from "$lib/interfaces";

export class LoginApi implements ILoginApi {
  private readonly __backendClient: IBackendClient;
  
  public email = $state<string>("");
  public password = $state<string>("");

  constructor(backendClient: IBackendClient) {
    this.__backendClient = backendClient;
  }

  public callApi = async (): Promise<string> => {
    await this.__backendClient.pb.collection("users").authWithPassword(this.email, this.password);
    return this.__backendClient.pb.authStore.exportToCookie({ httpOnly: false });
  };
}