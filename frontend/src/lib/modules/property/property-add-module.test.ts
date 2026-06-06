import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IPropertyService } from '$lib/services/property/property-service';
import { PropertyAddModule } from './property-add-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

const makeSvc = (overrides: Partial<IPropertyService> = {}) =>
  ({
    createPropertyWithOwner: vi.fn().mockResolvedValue({ id: 'pr1', address: '1 Main St' }),
    ...overrides,
  }) as unknown as IPropertyService;

describe('PropertyAddModule', () => {
  it('isAdd returns true', () => {
    const module = new PropertyAddModule(makeSvc(), 'person1', vi.fn());
    expect(module.isAdd).toBe(true);
  });

  it('callBackAction calls the provided callback', () => {
    const back = vi.fn();
    const module = new PropertyAddModule(makeSvc(), 'person1', back);
    module.callBackAction();
    expect(back).toHaveBeenCalledOnce();
  });

  it('resolves propertyAsync to an empty property model', async () => {
    const module = new PropertyAddModule(makeSvc(), 'person1', vi.fn());
    const property = await module.propertyAsync;
    expect(property).toMatchObject({ id: '', address: '' });
  });

  it('getSavePropertyAction calls createPropertyWithOwner with currentPersonId and navigates back', async () => {
    const back = vi.fn();
    const svc = makeSvc();
    const module = new PropertyAddModule(svc, 'person1', back);
    await module.getSavePropertyAction()({ id: '', address: '1 Main St' });
    expect(svc.createPropertyWithOwner).toHaveBeenCalledWith('1 Main St', 'person1');
    expect(back).toHaveBeenCalledOnce();
  });

  it('getDeletePropertyAction returns null', () => {
    const module = new PropertyAddModule(makeSvc(), 'person1', vi.fn());
    expect(module.getDeletePropertyAction()).toBeNull();
  });
});
