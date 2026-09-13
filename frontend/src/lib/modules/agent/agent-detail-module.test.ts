import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IAgentService } from '$lib/services/agent/agent-service';
import { AgentDetailModule } from './agent-detail-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

const mockAgent = { id: 'ag1', person: 'p1', cobrand: 'c1' };
const mockPropertyAgents = [{ id: 'pa1', agent: 'ag1', property: 'pr1' }];

const makeAgentService = (overrides: Partial<IAgentService> = {}) =>
  ({
    getAgentById: vi.fn().mockResolvedValue(mockAgent),
    getPropertyAgentsForAgent: vi.fn().mockResolvedValue(mockPropertyAgents),
    deleteAgent: vi.fn().mockResolvedValue(undefined),
    addPropertyAgent: vi.fn().mockResolvedValue({ id: 'pa2', agent: 'ag1', property: 'pr2' }),
    removePropertyAgent: vi.fn().mockResolvedValue(undefined),
    ...overrides,
  }) as unknown as IAgentService;

describe('AgentDetailModule', () => {
  it('isAdd returns false', () => {
    const module = new AgentDetailModule(makeAgentService(), 'ag1', vi.fn());
    expect(module.isAdd).toBe(false);
  });

  it('callBackAction calls the provided callback', () => {
    const backAction = vi.fn();
    const module = new AgentDetailModule(makeAgentService(), 'ag1', backAction);
    module.callBackAction();
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('resolves agent from service', async () => {
    const module = new AgentDetailModule(makeAgentService(), 'ag1', vi.fn());
    await expect(module.agentAsync).resolves.toEqual(mockAgent);
  });

  it('resolves propertyAgents from service', async () => {
    const module = new AgentDetailModule(makeAgentService(), 'ag1', vi.fn());
    await expect(module.propertyAgentsAsync).resolves.toEqual(mockPropertyAgents);
  });

  it('getSaveAgentAction returns null', () => {
    const module = new AgentDetailModule(makeAgentService(), 'ag1', vi.fn());
    expect(module.getSaveAgentAction()).toBeNull();
  });

  it('getDeleteAgentAction calls service and navigates back', async () => {
    const backAction = vi.fn();
    const agentService = makeAgentService();
    const module = new AgentDetailModule(agentService, 'ag1', backAction);
    await module.getDeleteAgentAction()!('ag1');
    expect(agentService.deleteAgent).toHaveBeenCalledWith('ag1');
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('getAddPropertyAgentAction calls service and refreshes propertyAgents', async () => {
    const agentService = makeAgentService();
    const module = new AgentDetailModule(agentService, 'ag1', vi.fn());
    await module.getAddPropertyAgentAction()('ag1', 'pr2');
    expect(agentService.addPropertyAgent).toHaveBeenCalledWith('ag1', 'pr2');
    expect(agentService.getPropertyAgentsForAgent).toHaveBeenCalledTimes(2);
  });

  it('getRemovePropertyAgentAction calls service and refreshes propertyAgents', async () => {
    const agentService = makeAgentService();
    const module = new AgentDetailModule(agentService, 'ag1', vi.fn());
    await module.getRemovePropertyAgentAction()('pa1');
    expect(agentService.removePropertyAgent).toHaveBeenCalledWith('pa1');
    expect(agentService.getPropertyAgentsForAgent).toHaveBeenCalledTimes(2);
  });
});
