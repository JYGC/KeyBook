import PocketBase from 'pocketbase';
import type { IPropertyAgentModel } from '$lib/models/agent-models';

export interface IPropertyAgentRepository {
  getPropertyAgentsByAgentId(agentId: string): Promise<IPropertyAgentModel[]>;
  getPropertyAgentsByPropertyId(propertyId: string): Promise<IPropertyAgentModel[]>;
  createPropertyAgent(agentId: string, propertyId: string): Promise<IPropertyAgentModel>;
  deletePropertyAgent(id: string): Promise<void>;
}

export class PropertyAgentRepository implements IPropertyAgentRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getPropertyAgentsByAgentId(agentId: string): Promise<IPropertyAgentModel[]> {
    const allPropertyAgents = await this.backendClient.collection('propertyAgents').getFullList<IPropertyAgentModel>();
    return allPropertyAgents.filter((propertyAgent) => propertyAgent.agent === agentId);
  }

  async getPropertyAgentsByPropertyId(propertyId: string): Promise<IPropertyAgentModel[]> {
    const allPropertyAgents = await this.backendClient.collection('propertyAgents').getFullList<IPropertyAgentModel>();
    return allPropertyAgents.filter((propertyAgent) => propertyAgent.property === propertyId);
  }

  async createPropertyAgent(agentId: string, propertyId: string): Promise<IPropertyAgentModel> {
    return await this.backendClient.collection('propertyAgents').create<IPropertyAgentModel>({
      agent: agentId,
      property: propertyId,
    });
  }

  async deletePropertyAgent(id: string): Promise<void> {
    await this.backendClient.collection('propertyAgents').delete(id);
  }
}
