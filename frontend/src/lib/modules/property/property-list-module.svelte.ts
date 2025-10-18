import PocketBase from "pocketbase";
import type { IPropertyListModule } from "$lib/modules/interfaces";
import type { IPropertyListItemModel } from "$lib/models/property-models";

export class PropertyListModule implements IPropertyListModule {
  private readonly __backendClient: PocketBase;

  public propertyListAsync = $derived.by<Promise<IPropertyListItemModel[]>>(async () => {
    try {
      if (this.__backendClient.authStore.record === null) {
        throw new Error("Cannot find loggedInUser.");
      }
      const items = await this.__backendClient.collection("properties").getFullList<IPropertyListItemModel>({
        filter: `owners.id ?~ "${this.__backendClient.authStore.record.id}"`,
        fields: "id,address"
      });
      return items;
    } catch (ex) {
      alert(ex);
      return [];
    }
  });

  constructor(backendClient: PocketBase) {
    this.__backendClient = backendClient;
  }
}