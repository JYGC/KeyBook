import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IBackendClient, IPersonListModule } from "$lib/interfaces";
import type { IPersonListItemModel } from "$lib/models/person-models";

export class PersonListModule implements IPersonListModule {
  private readonly __backendClient: IBackendClient;
  private readonly __propertyContext: PropertyContext;

  public personListAsync = $derived.by<Promise<IPersonListItemModel[]>>(async () => {
		try {
			const items = await this.__backendClient.pb.collection("personlistview").getFullList<IPersonListItemModel>({
				filter: `propertyid = "${this.__propertyContext.selectedPropertyId}"`,
				fields: "id,personid,personname,persontype,devicejsons",
			});
			return items;
		} catch (ex) {
			alert(ex);
			return [];
		}
	});

  constructor(backendClient: IBackendClient, propertyContext: PropertyContext) {
    this.__backendClient = backendClient;
    this.__propertyContext = propertyContext;
  }
}