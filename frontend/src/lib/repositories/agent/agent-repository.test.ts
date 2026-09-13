import PocketBase from 'pocketbase';
import { describe, it, expect, beforeAll } from 'vitest';
import { AgentRepository } from './agent-repository';
import { PropertyAgentRepository } from './property-agent-repository';

const PB_URL = 'http://192.168.8.144:8090';
const backendClient = new PocketBase(PB_URL);

beforeAll(async () => {
  // PocketBase Go v0.22 uses /api/admins (not _superusers which is v0.23+).
  const adminAuthResponse = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ identity: 'casperchen91@hotmail.com', password: 'w3m#@tpth100' }),
  });
  const adminAuthResult = await adminAuthResponse.json();
  backendClient.authStore.save(adminAuthResult.token, adminAuthResult.admin);
});

describe('AgentRepository', () => {
  it('creates, reads, and deletes an agent', async () => {
    const agentRepository = new AgentRepository(backendClient);

    const cobrand = await backendClient.collection('cobrands').create<{ id: string }>({ name: 'Agent Repo Test Cobrand' });
    const person = await backendClient
      .collection('persons')
      .create<{ id: string }>({ name: 'Agent Repo Test Person', DOB: '1990-01-01 00:00:00.000Z' });

    const created = await agentRepository.createAgent(person.id, cobrand.id);
    expect(created.person).toBe(person.id);
    expect(created.cobrand).toBe(cobrand.id);
    expect(created.id).toBeTruthy();

    const fetched = await agentRepository.getAgentById(created.id);
    expect(fetched.id).toBe(created.id);
    expect(fetched.person).toBe(person.id);

    const allAgents = await agentRepository.getAllAgents();
    expect(allAgents.some((agent) => agent.id === created.id)).toBe(true);

    await agentRepository.deleteAgent(created.id);
    await backendClient.collection('persons').delete(person.id);
    await backendClient.collection('cobrands').delete(cobrand.id);
  });
});

describe('PropertyAgentRepository', () => {
  it('creates, queries, and deletes a property agent', async () => {
    const agentRepository = new AgentRepository(backendClient);
    const propertyAgentRepository = new PropertyAgentRepository(backendClient);

    const cobrand = await backendClient.collection('cobrands').create<{ id: string }>({ name: 'PropAgent Test Cobrand' });
    const person = await backendClient
      .collection('persons')
      .create<{ id: string }>({ name: 'PropAgent Test Person', DOB: '1990-01-01 00:00:00.000Z' });
    const agent = await agentRepository.createAgent(person.id, cobrand.id);
    const property = await backendClient
      .collection('properties')
      .create<{ id: string }>({ address: '55 Agent Test Ave' });

    const created = await propertyAgentRepository.createPropertyAgent(agent.id, property.id);
    expect(created.agent).toBe(agent.id);
    expect(created.property).toBe(property.id);

    const byAgent = await propertyAgentRepository.getPropertyAgentsByAgentId(agent.id);
    expect(byAgent.some((propertyAgent) => propertyAgent.id === created.id)).toBe(true);

    const byProperty = await propertyAgentRepository.getPropertyAgentsByPropertyId(property.id);
    expect(byProperty.some((propertyAgent) => propertyAgent.id === created.id)).toBe(true);

    await propertyAgentRepository.deletePropertyAgent(created.id);
    await agentRepository.deleteAgent(agent.id);
    await backendClient.collection('properties').delete(property.id);
    await backendClient.collection('persons').delete(person.id);
    await backendClient.collection('cobrands').delete(cobrand.id);
  });
});
