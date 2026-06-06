import PocketBase from 'pocketbase';
import type { IPropertyAgentModel } from '$lib/models/agent-models';

export interface IPropertyAgentRepository {
  getByAgentId(agentId: string): Promise<IPropertyAgentModel[]>;
  getByPropertyId(propertyId: string): Promise<IPropertyAgentModel[]>;
  create(agentId: string, propertyId: string): Promise<IPropertyAgentModel>;
  delete(id: string): Promise<void>;
}

export class PropertyAgentRepository implements IPropertyAgentRepository {
  constructor(private readonly pb: PocketBase) {}

  async getByAgentId(agentId: string): Promise<IPropertyAgentModel[]> {
    const all = await this.pb.collection('propertyAgents').getFullList<IPropertyAgentModel>();
    return all.filter((pa) => pa.agent === agentId);
  }

  async getByPropertyId(propertyId: string): Promise<IPropertyAgentModel[]> {
    const all = await this.pb.collection('propertyAgents').getFullList<IPropertyAgentModel>();
    return all.filter((pa) => pa.property === propertyId);
  }

  async create(agentId: string, propertyId: string): Promise<IPropertyAgentModel> {
    return await this.pb.collection('propertyAgents').create<IPropertyAgentModel>({
      agent: agentId,
      property: propertyId,
    });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('propertyAgents').delete(id);
  }
}
