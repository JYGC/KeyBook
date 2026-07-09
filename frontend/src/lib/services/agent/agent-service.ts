import type { IAgentRepository } from '$lib/repositories/agent/agent-repository';
import type { IPropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';
import type { IAgentModel, IPropertyAgentModel } from '$lib/models/agent-models';

export interface IAgentService {
  validateDuplicatePropertyAgent(existing: IPropertyAgentModel[], agentId: string): void;
  getAllAgents(): Promise<IAgentModel[]>;
  getAgentById(id: string): Promise<IAgentModel>;
  createAgent(personId: string, cobrandId: string): Promise<IAgentModel>;
  deleteAgent(id: string): Promise<void>;
  getPropertyAgentsForAgent(agentId: string): Promise<IPropertyAgentModel[]>;
  addPropertyAgent(agentId: string, propertyId: string): Promise<IPropertyAgentModel>;
  removePropertyAgent(id: string): Promise<void>;
}

export class AgentService implements IAgentService {
  constructor(
    private readonly agentRepo: IAgentRepository,
    private readonly propertyAgentRepo: IPropertyAgentRepository,
  ) {}

  validateDuplicatePropertyAgent(existing: IPropertyAgentModel[], agentId: string): void {
    if (existing.some((pa) => pa.agent === agentId)) {
      throw new Error('agent is already assigned to this property');
    }
  }

  async getAllAgents(): Promise<IAgentModel[]> {
    return await this.agentRepo.getAll();
  }

  async getAgentById(id: string): Promise<IAgentModel> {
    return await this.agentRepo.getById(id);
  }

  async createAgent(personId: string, cobrandId: string): Promise<IAgentModel> {
    return await this.agentRepo.create(personId, cobrandId);
  }

  async deleteAgent(id: string): Promise<void> {
    await this.agentRepo.delete(id);
  }

  async getPropertyAgentsForAgent(agentId: string): Promise<IPropertyAgentModel[]> {
    return await this.propertyAgentRepo.getByAgentId(agentId);
  }

  async addPropertyAgent(agentId: string, propertyId: string): Promise<IPropertyAgentModel> {
    const existingForProperty = await this.propertyAgentRepo.getByPropertyId(propertyId);
    this.validateDuplicatePropertyAgent(existingForProperty, agentId);
    return await this.propertyAgentRepo.create(agentId, propertyId);
  }

  async removePropertyAgent(id: string): Promise<void> {
    await this.propertyAgentRepo.delete(id);
  }
}
