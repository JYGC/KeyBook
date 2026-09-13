import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IPropertyService } from '$lib/services/property/property-service';
import { PropertyAddModule } from './property-add-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

const makePropertyService = (overrides: Partial<IPropertyService> = {}) =>
  ({
    createPropertyWithOwner: vi.fn().mockResolvedValue({ id: 'pr1', address: '1 Main St' }),
    ...overrides,
  }) as unknown as IPropertyService;

describe('PropertyAddModule', () => {
  it('isAdd returns true', () => {
    const module = new PropertyAddModule(makePropertyService(), 'person1', vi.fn());
    expect(module.isAdd).toBe(true);
  });

  it('callBackAction calls the provided callback', () => {
    const backAction = vi.fn();
    const module = new PropertyAddModule(makePropertyService(), 'person1', backAction);
    module.callBackAction();
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('resolves propertyAsync to an empty property model', async () => {
    const module = new PropertyAddModule(makePropertyService(), 'person1', vi.fn());
    const property = await module.propertyAsync;
    expect(property).toMatchObject({ id: '', address: '' });
  });

  it('getSavePropertyAction calls createPropertyWithOwner with currentPersonId and navigates back', async () => {
    const backAction = vi.fn();
    const propertyService = makePropertyService();
    const module = new PropertyAddModule(propertyService, 'person1', backAction);
    await module.getSavePropertyAction()({ id: '', address: '1 Main St' });
    expect(propertyService.createPropertyWithOwner).toHaveBeenCalledWith('1 Main St', 'person1');
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('getDeletePropertyAction returns null', () => {
    const module = new PropertyAddModule(makePropertyService(), 'person1', vi.fn());
    expect(module.getDeletePropertyAction()).toBeNull();
  });
});
