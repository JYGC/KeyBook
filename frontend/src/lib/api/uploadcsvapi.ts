import type { IBackendClient, ICsvFileToObjectConverter, IUploadCsvApi } from "$lib/interfaces";

export class UploadCsvApi implements IUploadCsvApi {
  private readonly __csvFileToObjectConverter: ICsvFileToObjectConverter;
  private readonly __authManager: IBackendClient;
  
  constructor(
    csvFileToObjectConverter: ICsvFileToObjectConverter,
    authManager: IBackendClient
  ) {
    this.__csvFileToObjectConverter = csvFileToObjectConverter;
    this.__authManager = authManager;
  }

  public callApi = async () => {
    const csvContentObject = await this.__csvFileToObjectConverter.outputAsync;
    if (csvContentObject === null) {
      return;
    }
    await this.__authManager.pb.send("/api/device/importcsv", {
      method: "POST",
      body: csvContentObject,
    });
  };
}