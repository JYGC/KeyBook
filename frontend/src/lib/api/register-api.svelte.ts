import type { IBackendClient, IRegisterApi } from "$lib/interfaces";

export class RegisterApi implements IRegisterApi {
  private readonly __backendClient: IBackendClient;
  
  public name = $state<string>("");
  public email = $state<string>("");
  public password = $state<string>("");

  constructor(backendClient: IBackendClient) {
    this.__backendClient = backendClient;
  }

  private __error = $state<string>("");

  public callApi = async (): Promise<boolean> => {
    try {
      await this.__backendClient.pb.collection("users").create({
        name: this.name,
        password: this.password,
        passwordConfirm: this.password,
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