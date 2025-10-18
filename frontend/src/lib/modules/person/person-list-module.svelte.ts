import PocketBase from "pocketbase";
import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IPersonListModule } from "$lib/interfaces";
import type { IPersonListItemModel } from "$lib/models/person-device-models";

export class PersonListModule implements IPersonListModule {
  private readonly __backendClient: PocketBase;
  private readonly __propertyContext: PropertyContext;

  public personListAsync = $derived.by<Promise<IPersonListItemModel[]>>(async () => {
		try {
			const items = await this.__backendClient.collection("personlistview").getFullList<IPersonListItemModel>({
				filter: `propertyid = "${this.__propertyContext.selectedPropertyId}"`,
				fields: "id,personid,personname,persontype,holdingdevicejsons",
			})
			return items;
		} catch (ex) {
			alert(ex);
			return [];
		}
	});

  constructor(backendClient: PocketBase, propertyContext: PropertyContext) {
    this.__backendClient = backendClient;
    this.__propertyContext = propertyContext;
  }
}