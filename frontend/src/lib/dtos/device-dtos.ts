export interface IDeviceListItemDto {
  id: string;
  type: string;
  name: string;
  identifier: string;
}

export interface IEditDeviceDto {
  id: string;
  type: string;
  name: string;
  identifier: string;
  defunctreason: string;
  property: string;
}