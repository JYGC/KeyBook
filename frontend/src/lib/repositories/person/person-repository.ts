import PocketBase from 'pocketbase';
import type { IPersonModel } from '$lib/models/person-models';

export interface IPersonRepository {
  getAll(): Promise<IPersonModel[]>;
  getById(id: string): Promise<IPersonModel>;
  getByUserId(userId: string): Promise<IPersonModel | null>;
  create(name: string, dob: string): Promise<IPersonModel>;
  update(id: string, name: string, dob: string): Promise<IPersonModel>;
  delete(id: string): Promise<void>;
}

export class PersonRepository implements IPersonRepository {
  constructor(private readonly pb: PocketBase) {}

  async getAll(): Promise<IPersonModel[]> {
    return await this.pb.collection('persons').getFullList<IPersonModel>();
  }

  async getById(id: string): Promise<IPersonModel> {
    const all = await this.pb.collection('persons').getFullList<IPersonModel>();
    const person = all.find((p) => p.id === id);
    if (!person) throw new Error(`Person not found: ${id}`);
    return person;
  }

  async getByUserId(userId: string): Promise<IPersonModel | null> {
    const all = await this.pb.collection('persons').getFullList<IPersonModel>();
    return all.find((p) => p.user === userId) ?? null;
  }

  async create(name: string, dob: string): Promise<IPersonModel> {
    return await this.pb.collection('persons').create<IPersonModel>({ name, DOB: dob });
  }

  async update(id: string, name: string, dob: string): Promise<IPersonModel> {
    return await this.pb.collection('persons').update<IPersonModel>(id, { name, DOB: dob });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('persons').delete(id);
  }
}
