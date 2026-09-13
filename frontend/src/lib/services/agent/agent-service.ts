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
    private readonly agentRepository: IAgentRepository,
    private readonly propertyAgentRepository: IPropertyAgentRepository,
  ) {}

  validateDuplicatePropertyAgent(existing: IPropertyAgentModel[], agentId: string): void {
    if (existing.some((existingPropertyAgent) => existingPropertyAgent.agent === agentId)) {
      throw new Error('agent is already assigned to this property');
    }
  }

  async getAllAgents(): Promise<IAgentModel[]> {
    return await this.agentRepository.getAllAgents();
  }

  async getAgentById(id: string): Promise<IAgentModel> {
    return await this.agentRepository.getAgentById(id);
  }

  async createAgent(personId: string, cobrandId: string): Promise<IAgentModel> {
    return await this.agentRepository.createAgent(personId, cobrandId);
  }

  async deleteAgent(id: string): Promise<void> {
    await this.agentRepository.deleteAgent(id);
  }

  async getPropertyAgentsForAgent(agentId: string): Promise<IPropertyAgentModel[]> {
    return await this.propertyAgentRepository.getPropertyAgentsByAgentId(agentId);
  }

  async addPropertyAgent(agentId: string, propertyId: string): Promise<IPropertyAgentModel> {
    const existingForProperty = await this.propertyAgentRepository.getPropertyAgentsByPropertyId(propertyId);
    this.validateDuplicatePropertyAgent(existingForProperty, agentId);
    return await this.propertyAgentRepository.createPropertyAgent(agentId, propertyId);
  }

  async removePropertyAgent(id: string): Promise<void> {
    await this.propertyAgentRepository.deletePropertyAgent(id);
  }
}
