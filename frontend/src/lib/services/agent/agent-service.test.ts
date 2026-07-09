import { describe, it, expect, vi } from 'vitest';
import { AgentService } from './agent-service';
import type { IAgentRepository } from '$lib/repositories/agent/agent-repository';
import type { IPropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';

const mockAgentRepo: IAgentRepository = {
  getAll: vi.fn(),
  getById: vi.fn(),
  getByPersonId: vi.fn(),
  create: vi.fn(),
  delete: vi.fn(),
};

const mockPropertyAgentRepo: IPropertyAgentRepository = {
  getByAgentId: vi.fn(),
  getByPropertyId: vi.fn(),
  create: vi.fn(),
  delete: vi.fn(),
};

describe('AgentService.validateDuplicatePropertyAgent', () => {
  const svc = new AgentService(mockAgentRepo, mockPropertyAgentRepo);

  it('passes when agent is not already assigned to the property', () => {
    const existing = [
      { id: 'pa1', agent: 'a2', property: 'p1' },
      { id: 'pa2', agent: 'a3', property: 'p1' },
    ];
    expect(() => svc.validateDuplicatePropertyAgent(existing, 'a1')).not.toThrow();
  });

  it('passes for empty list', () => {
    expect(() => svc.validateDuplicatePropertyAgent([], 'a1')).not.toThrow();
  });

  it('throws when agent is already assigned to the property', () => {
    const existing = [
      { id: 'pa1', agent: 'a1', property: 'p1' },
    ];
    expect(() => svc.validateDuplicatePropertyAgent(existing, 'a1')).toThrow(
      'agent is already assigned to this property',
    );
  });
});
