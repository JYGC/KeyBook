import PocketBase from 'pocketbase';
import type { IPropertyModel } from '$lib/models/property-models';

export interface IPropertyRepository {
  getAllProperties(): Promise<IPropertyModel[]>;
  getPropertyById(id: string): Promise<IPropertyModel>;
  createProperty(address: string): Promise<IPropertyModel>;
  updateProperty(id: string, address: string): Promise<IPropertyModel>;
  deleteProperty(id: string): Promise<void>;
}

export class PropertyRepository implements IPropertyRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getAllProperties(): Promise<IPropertyModel[]> {
    return await this.backendClient.collection('properties').getFullList<IPropertyModel>();
  }

  async getPropertyById(id: string): Promise<IPropertyModel> {
    const allProperties = await this.backendClient.collection('properties').getFullList<IPropertyModel>();
    const property = allProperties.find((candidateProperty) => candidateProperty.id === id);
    if (!property) throw new Error(`Property not found: ${id}`);
    return property;
  }

  async createProperty(address: string): Promise<IPropertyModel> {
    return await this.backendClient.collection('properties').create<IPropertyModel>({ address });
  }

  async updateProperty(id: string, address: string): Promise<IPropertyModel> {
    return await this.backendClient.collection('properties').update<IPropertyModel>(id, { address });
  }

  async deleteProperty(id: string): Promise<void> {
    await this.backendClient.collection('properties').delete(id);
  }
}
