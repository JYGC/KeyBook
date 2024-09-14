import type { PropertyContext } from "$lib/contexts/property-context.svelte";
import type { IPersonDeviceModel } from "$lib/dtos/person-device-models";
import type { IPersonDeviceExpandPersonDevicePersonEditModel, IPersonIdNameTypeModel } from "$lib/dtos/person-models";
import type { IBackendClient, IDeviceEditorModule, IDeviceHolderEditorModule } from "$lib/interfaces";

export class DeviceHolderEditorModule implements IDeviceHolderEditorModule {
  private readonly __backendClient: IBackendClient;
  private readonly __deviceEditorModule: IDeviceEditorModule;
  private readonly __propertyContext: PropertyContext;

  private __getPersonDeviceExpandPersonDevicePersonAsync = async () => {
    try {
      const deviceId = (await this.__deviceEditorModule.deviceAsync)?.id;
      if (deviceId === undefined || deviceId === null) {
        return null;
      }
      return (await this.__backendClient.pb.collection("persondevices").getFirstListItem<IPersonDeviceExpandPersonDevicePersonEditModel>(`device = "${deviceId}"`, {
        expand: "person",
        fields: "id,expand.person.id,expand.person.name,expand.person.type",
      }));
    } catch (ex) {
      if (!(ex instanceof Error)) {
        alert("Error");
      } else if (ex.message !== "The requested resource wasn't found.") {
        alert(ex);
      }
      return null;
    }
  };

  public personDeviceExpandPersonDevicePersonAsync = $state<Promise<IPersonDeviceExpandPersonDevicePersonEditModel | null>>((async () => null)());

  public availablePersonsAsync = $derived.by<Promise<IPersonIdNameTypeModel[]>>(async () => {
    try {
      if (this.__backendClient.loggedInUser === null) {
        throw new Error("Cannot find loggedInUser.");
      }
      return (await this.__backendClient.pb.collection("persons").getList<IPersonIdNameTypeModel>(1, 50, {
        filter: `property.id = "${this.__propertyContext.selectedPropertyId}" && property.owners.id ?~ "${this.__backendClient.loggedInUser.id}"`,
        fields: "id,name,type"
      })).items;
    } catch {
      return [];
    }
  });

  public replaceDeviceHolderActionAsync = async () => {
    try {
      const deviceId = (await this.__deviceEditorModule.deviceAsync)?.id;
      if (deviceId === undefined) {
        return;
      }
      const personDeviceId = (await this.personDeviceExpandPersonDevicePersonAsync)?.id;
      if (personDeviceId !== undefined) {
        await this.__backendClient.pb.collection("persondevices").delete(personDeviceId);
      }

      if (this.selectedDeviceHolderId !== "") {
        await this.__backendClient.pb.collection("persondevices").create<IPersonDeviceModel>({
          person: this.selectedDeviceHolderId,
          device: deviceId,
        });
      }

      this.personDeviceExpandPersonDevicePersonAsync = this.__getPersonDeviceExpandPersonDevicePersonAsync();
    } catch (ex) {
      alert(ex);
    }
  };

  public currentDeviceHolderNameAsync = $derived.by<Promise<string>>(async () => (await this.personDeviceExpandPersonDevicePersonAsync)?.expand.person.name ?? "");

  public selectedDeviceHolderId = $state<string>("");

  constructor(
    backendClient: IBackendClient,
    deviceEditorModule: IDeviceEditorModule,
    propertyContext: PropertyContext,
  ) {
    this.__backendClient = backendClient;
    this.__deviceEditorModule = deviceEditorModule;
    this.__propertyContext = propertyContext;

    this.personDeviceExpandPersonDevicePersonAsync = this.__getPersonDeviceExpandPersonDevicePersonAsync();

    this.personDeviceExpandPersonDevicePersonAsync
      .then((personDeviceExpandPersonDevicePerson: IPersonDeviceExpandPersonDevicePersonEditModel | null) => {
        this.selectedDeviceHolderId = personDeviceExpandPersonDevicePerson?.expand.person.id ?? "";
      });
  }
}