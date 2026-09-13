import { describe, it, expect, vi } from 'vitest';
import { AgentService } from './agent-service';
import type { IAgentRepository } from '$lib/repositories/agent/agent-repository';
import type { IPropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';

const mockAgentRepository: IAgentRepository = {
  getAllAgents: vi.fn(),
  getAgentById: vi.fn(),
  getAgentsByPersonId: vi.fn(),
  createAgent: vi.fn(),
  deleteAgent: vi.fn(),
};

const mockPropertyAgentRepository: IPropertyAgentRepository = {
  getPropertyAgentsByAgentId: vi.fn(),
  getPropertyAgentsByPropertyId: vi.fn(),
  createPropertyAgent: vi.fn(),
  deletePropertyAgent: vi.fn(),
};

describe('AgentService.validateDuplicatePropertyAgent', () => {
  const agentService = new AgentService(mockAgentRepository, mockPropertyAgentRepository);

  it('passes when agent is not already assigned to the property', () => {
    const existing = [
      { id: 'pa1', agent: 'a2', property: 'p1' },
      { id: 'pa2', agent: 'a3', property: 'p1' },
    ];
    expect(() => agentService.validateDuplicatePropertyAgent(existing, 'a1')).not.toThrow();
  });

  it('passes for empty list', () => {
    expect(() => agentService.validateDuplicatePropertyAgent([], 'a1')).not.toThrow();
  });

  it('throws when agent is already assigned to the property', () => {
    const existing = [
      { id: 'pa1', agent: 'a1', property: 'p1' },
    ];
    expect(() => agentService.validateDuplicatePropertyAgent(existing, 'a1')).toThrow(
      'agent is already assigned to this property',
    );
  });
});
