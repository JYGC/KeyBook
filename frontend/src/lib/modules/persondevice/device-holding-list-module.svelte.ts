import type { PersonContext } from "$lib/contexts/person-context.svelte";
import type { IBackendClient, IDeviceHoldingListModule } from "$lib/interfaces";
import type { IDeviceHeldExpandDeviceModel, IDeviceListItemModel } from "$lib/models/person-device-models";

export class DeviceHoldingListModule implements IDeviceHoldingListModule {
  private readonly __backendClient: IBackendClient;
  private readonly __personContext: PersonContext;
  
  public deviceHoldingListOfPersonAsync = $derived.by<Promise<IDeviceListItemModel[]>>(async () => {
    try {
      if (this.__backendClient.loggedInUser === null) {
        throw new Error("Cannot find loggedInUser.");
      }
      const items = await this.__backendClient.pb.collection("persondevices").getFullList<IDeviceHeldExpandDeviceModel>({
        filter: `person.id = "${this.__personContext.selectedPersonId}"`,
        expand: "device",
        fields: "id,expand.device.id,expand.device.name,expand.device.identifier,expand.device.type",
      });
      return items.map(dh => dh.expand.device);
    } catch {
      return [];
    }
  });

  constructor(
    backendClient: IBackendClient,
    personContext: PersonContext
  ) {
    this.__backendClient = backendClient;
    this.__personContext = personContext
  }
}