import PocketBase from 'pocketbase';
import { describe, it, expect, beforeAll } from 'vitest';
import { AgentRepository } from './agent-repository';
import { PropertyAgentRepository } from './property-agent-repository';

const PB_URL = 'http://192.168.8.144:8090';
const pb = new PocketBase(PB_URL);

beforeAll(async () => {
  // PocketBase Go v0.22 uses /api/admins (not _superusers which is v0.23+).
  const res = await fetch(`${PB_URL}/api/admins/auth-with-password`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ identity: 'casperchen91@hotmail.com', password: 'w3m#@tpth100' }),
  });
  const data = await res.json();
  pb.authStore.save(data.token, data.admin);
});

describe('AgentRepository', () => {
  it('creates, reads, and deletes an agent', async () => {
    const repo = new AgentRepository(pb);

    const cobrand = await pb.collection('cobrands').create<{ id: string }>({ name: 'Agent Repo Test Cobrand' });
    const person = await pb
      .collection('persons')
      .create<{ id: string }>({ name: 'Agent Repo Test Person', DOB: '1990-01-01 00:00:00.000Z' });

    const created = await repo.create(person.id, cobrand.id);
    expect(created.person).toBe(person.id);
    expect(created.cobrand).toBe(cobrand.id);
    expect(created.id).toBeTruthy();

    const fetched = await repo.getById(created.id);
    expect(fetched.id).toBe(created.id);
    expect(fetched.person).toBe(person.id);

    const all = await repo.getAll();
    expect(all.some((a) => a.id === created.id)).toBe(true);

    await repo.delete(created.id);
    await pb.collection('persons').delete(person.id);
    await pb.collection('cobrands').delete(cobrand.id);
  });
});

describe('PropertyAgentRepository', () => {
  it('creates, queries, and deletes a property agent', async () => {
    const agentRepo = new AgentRepository(pb);
    const propertyAgentRepo = new PropertyAgentRepository(pb);

    const cobrand = await pb.collection('cobrands').create<{ id: string }>({ name: 'PropAgent Test Cobrand' });
    const person = await pb
      .collection('persons')
      .create<{ id: string }>({ name: 'PropAgent Test Person', DOB: '1990-01-01 00:00:00.000Z' });
    const agent = await agentRepo.create(person.id, cobrand.id);
    const property = await pb
      .collection('properties')
      .create<{ id: string }>({ address: '55 Agent Test Ave' });

    const created = await propertyAgentRepo.create(agent.id, property.id);
    expect(created.agent).toBe(agent.id);
    expect(created.property).toBe(property.id);

    const byAgent = await propertyAgentRepo.getByAgentId(agent.id);
    expect(byAgent.some((pa) => pa.id === created.id)).toBe(true);

    const byProperty = await propertyAgentRepo.getByPropertyId(property.id);
    expect(byProperty.some((pa) => pa.id === created.id)).toBe(true);

    await propertyAgentRepo.delete(created.id);
    await agentRepo.delete(agent.id);
    await pb.collection('properties').delete(property.id);
    await pb.collection('persons').delete(person.id);
    await pb.collection('cobrands').delete(cobrand.id);
  });
});
