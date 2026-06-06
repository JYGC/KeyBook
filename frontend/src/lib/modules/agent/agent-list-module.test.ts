import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IAgentService } from '$lib/services/agent/agent-service';
import { AgentListModule } from './agent-list-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

describe('AgentListModule', () => {
  it('returns agents from service', async () => {
    const agents = [
      { id: '1', person: 'p1', cobrand: 'c1' },
      { id: '2', person: 'p2', cobrand: 'c1' },
    ];
    const svc = { getAllAgents: vi.fn().mockResolvedValue(agents) } as unknown as IAgentService;

    const module = new AgentListModule(svc);
    await expect(module.agentListAsync).resolves.toEqual(agents);
    expect(svc.getAllAgents).toHaveBeenCalledOnce();
  });

  it('returns empty array when service throws', async () => {
    const svc = {
      getAllAgents: vi.fn().mockRejectedValue(new Error('network error')),
    } as unknown as IAgentService;

    const module = new AgentListModule(svc);
    await expect(module.agentListAsync).resolves.toEqual([]);
  });
});
