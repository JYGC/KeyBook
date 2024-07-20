import type { AddDeviceAndHistoriesDTO, AddPropertyDeviceAndHistoriesDTO } from "$lib/dtos";
import type { ICsvFileToObjectConverter } from "$lib/interfaces";

export class CsvFileToObjectConverter implements ICsvFileToObjectConverter {  
  private convertCsvFileToObject = async (input: FileList | null): Promise<AddPropertyDeviceAndHistoriesDTO | null> => {
    if (input === null) {
      return null;
    }
    const uploadedFile = input[0];
    const csvText = await uploadedFile.text(); 

    const lines = csvText.split('\n');
  
    const nonEmptyLines: string[] = [];
  
    const regex = /[a-zA-Z0-9]/;
    lines.forEach((line) => {
      if (regex.test(line)) {
        nonEmptyLines.push(line.trim());
      }
    });
  
    const propertyLineSplit = nonEmptyLines[0].split('"');
  
    const propertyAddress = propertyLineSplit[1];
  
    const subHeadingSplit = nonEmptyLines[1].split(',');
  
    const devicesPersonDevicesAndHistories = [];
    for (let c = 2; c < nonEmptyLines.length; c++) {
      const deviceLineSplit = nonEmptyLines[c].split(',');
  
      let deviceProcurementDate: Date | null = null;
  
      const device: AddDeviceAndHistoriesDTO = {
        Name: deviceLineSplit[0],
        Identifier: deviceLineSplit[1],
        Type: deviceLineSplit[2],
        DefunctReason: 'None',
        CurrentHolder: null,
        DeviceHistories: [],
        PersonDeviceHistories: []
      };
  
      let lastNonEmptyRow = -1;
      for (let r = 3; r < deviceLineSplit.length; r++) {
        const deviceHolder = deviceLineSplit[r].trim();
        if (deviceHolder.length > 0) {
          lastNonEmptyRow = r;
          if (deviceProcurementDate === null) {
            deviceProcurementDate = new Date(subHeadingSplit[r]);
            device.DeviceHistories.push({
              dateSpecified: deviceProcurementDate,
              actionDescription: 'procured'
            });
          }
          if (deviceHolder === '!!LOST!!') {
            device.PersonDeviceHistories.push({
              DeviceHolder: null,
              DateSpecified: new Date(subHeadingSplit[r]),
              ActionDescription: 'device lost'
            });
          } else if (deviceHolder === '!!STOLEN BY TENANT!!') {
            device.PersonDeviceHistories.push({
              DeviceHolder: null,
              DateSpecified: new Date(subHeadingSplit[r]),
              ActionDescription: 'device stolen'
            });
          } else if (deviceHolder === '!!NOT WORKING!!') {
            device.PersonDeviceHistories.push({
              DeviceHolder: null,
              DateSpecified: new Date(subHeadingSplit[r]),
              ActionDescription: 'device damaged'
            });
          } else if (deviceHolder !== 'Storage') {
            device.PersonDeviceHistories.push({
              DeviceHolder: deviceHolder,
              DateSpecified: new Date(subHeadingSplit[r]),
              ActionDescription: 'device held by'
            });
          }
        }
      }
  
      if (lastNonEmptyRow < deviceLineSplit.length - 1) {
        const lastDeviceHistory: {
          dateSpecified: Date,
          actionDescription: string | null
        } = {
          dateSpecified:  new Date(subHeadingSplit[lastNonEmptyRow]),
          actionDescription: null
        };
  
        const lastDeviceHolder = deviceLineSplit[lastNonEmptyRow].trim();
        if (lastDeviceHolder === '!!LOST!!') {
          device.DefunctReason = 'Lost';
          lastDeviceHistory.actionDescription = 'device lost';
        } else if (lastDeviceHolder === '!!STOLEN BY TENANT!!') {
          device.DefunctReason = 'Stolen';
          lastDeviceHistory.actionDescription = 'device stolen';
        } else if (lastDeviceHolder === '!!NOT WORKING!!') {
          device.DefunctReason = 'Damaged';
          lastDeviceHistory.actionDescription = 'device damaged';
        } else {
          device.DefunctReason = 'Retired';
          lastDeviceHistory.actionDescription = 'device retired';
        }
  
        device.DeviceHistories.push(lastDeviceHistory);
      }
  
      const currentHolderRegistered = deviceLineSplit[deviceLineSplit.length - 1];
  
      device.CurrentHolder = currentHolderRegistered.length === 0 || currentHolderRegistered === 'Storage' || currentHolderRegistered === '!!LOST!!' ? null : currentHolderRegistered;
  
      devicesPersonDevicesAndHistories.push(device);
    }

    return {
      PropertyAddress: propertyAddress,
      DevicesPersonDevicesAndHistories: devicesPersonDevicesAndHistories
    }
  };
  
  public input = $state<FileList | null>(null);
  public outputAsync = $derived<Promise<AddPropertyDeviceAndHistoriesDTO | null>>(this.convertCsvFileToObject(this.input));

  constructor() { }
}