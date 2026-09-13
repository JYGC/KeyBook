import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IPersonService } from '$lib/services/person/person-service';
import { PersonDetailModule } from './person-detail-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

const mockPerson = { id: 'p1', name: 'Alice', DOB: '1990-01-01', user: 'u1', profileImage: '' };
const mockRoles = ['Owner', 'Tenant'];

const makePersonService = (overrides: Partial<IPersonService> = {}) =>
  ({
    getPersonById: vi.fn().mockResolvedValue(mockPerson),
    getRolesForPerson: vi.fn().mockResolvedValue(mockRoles),
    updatePerson: vi.fn().mockResolvedValue(mockPerson),
    deletePerson: vi.fn().mockResolvedValue(undefined),
    ...overrides,
  }) as unknown as IPersonService;

describe('PersonDetailModule', () => {
  it('isAdd returns false', () => {
    const module = new PersonDetailModule(makePersonService(), 'p1', vi.fn());
    expect(module.isAdd).toBe(false);
  });

  it('callBackAction calls the provided callback', () => {
    const backAction = vi.fn();
    const module = new PersonDetailModule(makePersonService(), 'p1', backAction);
    module.callBackAction();
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('resolves person from service', async () => {
    const module = new PersonDetailModule(makePersonService(), 'p1', vi.fn());
    await expect(module.personAsync).resolves.toEqual(mockPerson);
  });

  it('resolves roles from service', async () => {
    const module = new PersonDetailModule(makePersonService(), 'p1', vi.fn());
    await expect(module.rolesAsync).resolves.toEqual(mockRoles);
  });

  it('getSavePersonAction calls updatePerson and navigates back', async () => {
    const backAction = vi.fn();
    const personService = makePersonService();
    const module = new PersonDetailModule(personService, 'p1', backAction);
    await module.getSavePersonAction()(mockPerson);
    expect(personService.updatePerson).toHaveBeenCalledWith('p1', 'Alice', '1990-01-01');
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('getDeletePersonAction calls deletePerson and navigates back', async () => {
    const backAction = vi.fn();
    const personService = makePersonService();
    const module = new PersonDetailModule(personService, 'p1', backAction);
    await module.getDeletePersonAction()!('p1');
    expect(personService.deletePerson).toHaveBeenCalledWith('p1');
    expect(backAction).toHaveBeenCalledOnce();
  });
});
