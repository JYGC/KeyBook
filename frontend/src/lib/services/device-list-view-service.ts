import PocketBase, { RecordService } from "pocketbase";
import type { IDeviceListItemModel } from "$lib/models/person-device-models";
import type { IDeviceListViewService } from "./interfaces";

export class DeviceListViewService implements IDeviceListViewService {
  private readonly __deviceListViewCollection: RecordService<IDeviceListItemModel>;

  constructor(backendClient: PocketBase) {
    this.__deviceListViewCollection = backendClient.collection<IDeviceListItemModel>("devicelistview");
  }
  public getForDeviceListAsync = async (selectedPropertyId: string): Promise<IDeviceListItemModel[]> => {
    return await this.__deviceListViewCollection.getFullList({
      filter: `propertyid = "${selectedPropertyId}"`,
      fields: "id,deviceid,devicename,deviceidentifier,devicetype,personname",
    });
  };
}