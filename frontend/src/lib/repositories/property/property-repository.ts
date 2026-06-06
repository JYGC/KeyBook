import PocketBase from 'pocketbase';
import type { IPropertyModel } from '$lib/models/property-models';

export interface IPropertyRepository {
  getAll(): Promise<IPropertyModel[]>;
  getById(id: string): Promise<IPropertyModel>;
  create(address: string): Promise<IPropertyModel>;
  update(id: string, address: string): Promise<IPropertyModel>;
  delete(id: string): Promise<void>;
}

export class PropertyRepository implements IPropertyRepository {
  constructor(private readonly pb: PocketBase) {}

  async getAll(): Promise<IPropertyModel[]> {
    return await this.pb.collection('properties').getFullList<IPropertyModel>();
  }

  async getById(id: string): Promise<IPropertyModel> {
    const all = await this.pb.collection('properties').getFullList<IPropertyModel>();
    const property = all.find((p) => p.id === id);
    if (!property) throw new Error(`Property not found: ${id}`);
    return property;
  }

  async create(address: string): Promise<IPropertyModel> {
    return await this.pb.collection('properties').create<IPropertyModel>({ address });
  }

  async update(id: string, address: string): Promise<IPropertyModel> {
    return await this.pb.collection('properties').update<IPropertyModel>(id, { address });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('properties').delete(id);
  }
}
