export interface IQueryParameters {
  page: number;
  perPage: number;
  parameters: {
    expand: string | null,
    fields: string | null,
    filter: string | null
  };
}

export interface AddDeviceAndHistoriesDTO {
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

export interface AddPropertyDeviceAndHistoriesDTO {
  PropertyAddress: string;
  DevicesPersonDevicesAndHistories: AddDeviceAndHistoriesDTO[];
}

export interface DeviceIdNameIdentifierAndTypeDTO {
  Name: string,
  Identifier: string,
  Type: string,
}