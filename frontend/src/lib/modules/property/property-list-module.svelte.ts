import type { IBackendClient, IPropertyListModule } from "$lib/interfaces";
import type { IPropertyListItemModel } from "$lib/models/property-models";

export class PropertyListModule implements IPropertyListModule {
  private readonly __backendClient: IBackendClient;

  public propertyListAsync = $derived.by<Promise<IPropertyListItemModel[]>>(async () => {
    try {
      if (this.__backendClient.loggedInUser === null) {
        throw new Error("Cannot find loggedInUser.");
      }
      const response = await this.__backendClient.pb.collection("properties").getList<IPropertyListItemModel>(1, 50, {
        filter: `owners.id ?~ "${this.__backendClient.loggedInUser.id}"`,
        fields: "id,address"
      });
      return response.items;
    } catch (ex) {
      alert(ex);
      return [];
    }
  });

  constructor(backendClient: IBackendClient) {
    this.__backendClient = backendClient;
  }
}