import PocketBase from 'pocketbase';
import type { IAgentModel } from '$lib/models/agent-models';

export interface IAgentRepository {
  getAll(): Promise<IAgentModel[]>;
  getById(id: string): Promise<IAgentModel>;
  create(personId: string, cobrandId: string): Promise<IAgentModel>;
  delete(id: string): Promise<void>;
}

export class AgentRepository implements IAgentRepository {
  constructor(private readonly pb: PocketBase) {}

  async getAll(): Promise<IAgentModel[]> {
    return await this.pb.collection('agents').getFullList<IAgentModel>();
  }

  async getById(id: string): Promise<IAgentModel> {
    // PocketBase v0.22: combining id= filter with a back-relation listRule
    // produces incorrect SQL. Fetch all accessible records and match client-side.
    const all = await this.pb.collection('agents').getFullList<IAgentModel>();
    const agent = all.find((a) => a.id === id);
    if (!agent) throw new Error(`Agent not found: ${id}`);
    return agent;
  }

  async create(personId: string, cobrandId: string): Promise<IAgentModel> {
    return await this.pb.collection('agents').create<IAgentModel>({
      person: personId,
      cobrand: cobrandId,
    });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('agents').delete(id);
  }
}
