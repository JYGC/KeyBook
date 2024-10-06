import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IBackendClient, IPersonListModule } from "$lib/interfaces";
import type { IPersonIdNameTypeModel } from "$lib/models/person-models";

export class PersonListModule implements IPersonListModule {
  private readonly __backendClient: IBackendClient;
  private readonly __propertyContext: PropertyContext;

  public personListAsync = $derived.by<Promise<IPersonIdNameTypeModel[]>>(async () => {
		try {
			const response = await this.__backendClient.pb.collection("persons").getList<IPersonIdNameTypeModel>(
				1,
				50,
				{
					filter: `property.id = "${this.__propertyContext.selectedPropertyId}"`,
					fields: "id,type,name",
				}
			);
			return response.items;
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