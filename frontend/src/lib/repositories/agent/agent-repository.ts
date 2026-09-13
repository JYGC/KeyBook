import PocketBase from 'pocketbase';
import type { IAgentModel } from '$lib/models/agent-models';

export interface IAgentRepository {
  getAllAgents(): Promise<IAgentModel[]>;
  getAgentById(id: string): Promise<IAgentModel>;
  getAgentsByPersonId(personId: string): Promise<IAgentModel[]>;
  createAgent(personId: string, cobrandId: string): Promise<IAgentModel>;
  deleteAgent(id: string): Promise<void>;
}

export class AgentRepository implements IAgentRepository {
  constructor(private readonly backendClient: PocketBase) {}

  async getAllAgents(): Promise<IAgentModel[]> {
    return await this.backendClient.collection('agents').getFullList<IAgentModel>();
  }

  async getAgentById(id: string): Promise<IAgentModel> {
    // PocketBase v0.22: combining id= filter with a back-relation listRule
    // produces incorrect SQL. Fetch all accessible records and match client-side.
    const allAgents = await this.backendClient.collection('agents').getFullList<IAgentModel>();
    const agent = allAgents.find((candidateAgent) => candidateAgent.id === id);
    if (!agent) throw new Error(`Agent not found: ${id}`);
    return agent;
  }

  async getAgentsByPersonId(personId: string): Promise<IAgentModel[]> {
    const allAgents = await this.backendClient.collection('agents').getFullList<IAgentModel>();
    return allAgents.filter((agent) => agent.person === personId);
  }

  async createAgent(personId: string, cobrandId: string): Promise<IAgentModel> {
    return await this.backendClient.collection('agents').create<IAgentModel>({
      person: personId,
      cobrand: cobrandId,
    });
  }

  async deleteAgent(id: string): Promise<void> {
    await this.backendClient.collection('agents').delete(id);
  }
}
