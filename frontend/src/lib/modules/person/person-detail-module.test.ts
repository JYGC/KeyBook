import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IPersonService } from '$lib/services/person/person-service';
import { PersonDetailModule } from './person-detail-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

const mockPerson = { id: 'p1', name: 'Alice', DOB: '1990-01-01', user: 'u1', profileImage: '' };
const mockRoles = ['Owner', 'Tenant'];

const makeSvc = (overrides: Partial<IPersonService> = {}) =>
  ({
    getPersonById: vi.fn().mockResolvedValue(mockPerson),
    getRolesForPerson: vi.fn().mockResolvedValue(mockRoles),
    updatePerson: vi.fn().mockResolvedValue(mockPerson),
    deletePerson: vi.fn().mockResolvedValue(undefined),
    ...overrides,
  }) as unknown as IPersonService;

describe('PersonDetailModule', () => {
  it('isAdd returns false', () => {
    const module = new PersonDetailModule(makeSvc(), 'p1', vi.fn());
    expect(module.isAdd).toBe(false);
  });

  it('callBackAction calls the provided callback', () => {
    const back = vi.fn();
    const module = new PersonDetailModule(makeSvc(), 'p1', back);
    module.callBackAction();
    expect(back).toHaveBeenCalledOnce();
  });

  it('resolves person from service', async () => {
    const module = new PersonDetailModule(makeSvc(), 'p1', vi.fn());
    await expect(module.personAsync).resolves.toEqual(mockPerson);
  });

  it('resolves roles from service', async () => {
    const module = new PersonDetailModule(makeSvc(), 'p1', vi.fn());
    await expect(module.rolesAsync).resolves.toEqual(mockRoles);
  });

  it('getSavePersonAction calls updatePerson and navigates back', async () => {
    const back = vi.fn();
    const svc = makeSvc();
    const module = new PersonDetailModule(svc, 'p1', back);
    await module.getSavePersonAction()(mockPerson);
    expect(svc.updatePerson).toHaveBeenCalledWith('p1', 'Alice', '1990-01-01');
    expect(back).toHaveBeenCalledOnce();
  });

  it('getDeletePersonAction calls deletePerson and navigates back', async () => {
    const back = vi.fn();
    const svc = makeSvc();
    const module = new PersonDetailModule(svc, 'p1', back);
    await module.getDeletePersonAction()!('p1');
    expect(svc.deletePerson).toHaveBeenCalledWith('p1');
    expect(back).toHaveBeenCalledOnce();
  });
});
