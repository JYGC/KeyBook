export interface AddDeviceAndHistoriesModel {
  Name: string;
  Identifier: string;
  Type: string;
  DefunctReason: string;
  CurrentHolder: string | null;
  DeviceHistories: object[];
  PersonDeviceHistories: {
    DeviceHolder: string | null;
    DateSpecified: Date;
    ActionDescription: string;
  }[]
}

export interface AddPropertyDeviceAndHistoriesModel {
  PropertyAddress: string;
  DevicesPersonDevicesAndHistories: AddDeviceAndHistoriesModel[];
}

export interface DeviceIdNameIdentifierAndTypeModel {
  Name: string,
  Identifier: string,
  Type: string,
}