import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IPersonService } from '$lib/services/person/person-service';
import { PersonListModule } from './person-list-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

describe('PersonListModule', () => {
  it('returns persons from service', async () => {
    const persons = [
      { id: '1', name: 'Alice', DOB: '1990-01-01', user: 'u1', profileImage: '' },
      { id: '2', name: 'Bob', DOB: '1985-05-15', user: 'u2', profileImage: '' },
    ];
    const personService = { getAllPersons: vi.fn().mockResolvedValue(persons) } as unknown as IPersonService;

    const module = new PersonListModule(personService);
    await expect(module.personListAsync).resolves.toEqual(persons);
    expect(personService.getAllPersons).toHaveBeenCalledOnce();
  });

  it('returns empty array when service throws', async () => {
    const personService = {
      getAllPersons: vi.fn().mockRejectedValue(new Error('network error')),
    } as unknown as IPersonService;

    const module = new PersonListModule(personService);
    await expect(module.personListAsync).resolves.toEqual([]);
  });
});
