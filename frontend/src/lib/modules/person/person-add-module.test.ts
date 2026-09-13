import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IPersonService } from '$lib/services/person/person-service';
import { PersonAddModule } from './person-add-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

const makePersonService = (overrides: Partial<IPersonService> = {}) =>
  ({
    createPerson: vi.fn().mockResolvedValue({ id: 'p1', name: 'Alice', DOB: '1990-01-01', user: '', profileImage: '' }),
    ...overrides,
  }) as unknown as IPersonService;

describe('PersonAddModule', () => {
  it('isAdd returns true', () => {
    const module = new PersonAddModule(makePersonService(), vi.fn());
    expect(module.isAdd).toBe(true);
  });

  it('callBackAction calls the provided callback', () => {
    const backAction = vi.fn();
    const module = new PersonAddModule(makePersonService(), backAction);
    module.callBackAction();
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('resolves personAsync to an empty person model', async () => {
    const module = new PersonAddModule(makePersonService(), vi.fn());
    const person = await module.personAsync;
    expect(person).toMatchObject({ id: '', name: '', DOB: '' });
  });

  it('getSavePersonAction calls createPerson and navigates back', async () => {
    const backAction = vi.fn();
    const personService = makePersonService();
    const module = new PersonAddModule(personService, backAction);
    await module.getSavePersonAction()({ id: '', name: 'Alice', DOB: '1990-01-01', user: '', profileImage: '' });
    expect(personService.createPerson).toHaveBeenCalledWith('Alice', '1990-01-01');
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('getDeletePersonAction returns null', () => {
    const module = new PersonAddModule(makePersonService(), vi.fn());
    expect(module.getDeletePersonAction()).toBeNull();
  });
});
