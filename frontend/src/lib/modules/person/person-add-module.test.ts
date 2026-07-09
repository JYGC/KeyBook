import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IPersonService } from '$lib/services/person/person-service';
import { PersonAddModule } from './person-add-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

const makeSvc = (overrides: Partial<IPersonService> = {}) =>
  ({
    createPerson: vi.fn().mockResolvedValue({ id: 'p1', name: 'Alice', DOB: '1990-01-01', user: '', profileImage: '' }),
    ...overrides,
  }) as unknown as IPersonService;

describe('PersonAddModule', () => {
  it('isAdd returns true', () => {
    const module = new PersonAddModule(makeSvc(), vi.fn());
    expect(module.isAdd).toBe(true);
  });

  it('callBackAction calls the provided callback', () => {
    const back = vi.fn();
    const module = new PersonAddModule(makeSvc(), back);
    module.callBackAction();
    expect(back).toHaveBeenCalledOnce();
  });

  it('resolves personAsync to an empty person model', async () => {
    const module = new PersonAddModule(makeSvc(), vi.fn());
    const person = await module.personAsync;
    expect(person).toMatchObject({ id: '', name: '', DOB: '' });
  });

  it('getSavePersonAction calls createPerson and navigates back', async () => {
    const back = vi.fn();
    const svc = makeSvc();
    const module = new PersonAddModule(svc, back);
    await module.getSavePersonAction()({ id: '', name: 'Alice', DOB: '1990-01-01', user: '', profileImage: '' });
    expect(svc.createPerson).toHaveBeenCalledWith('Alice', '1990-01-01');
    expect(back).toHaveBeenCalledOnce();
  });

  it('getDeletePersonAction returns null', () => {
    const module = new PersonAddModule(makeSvc(), vi.fn());
    expect(module.getDeletePersonAction()).toBeNull();
  });
});
