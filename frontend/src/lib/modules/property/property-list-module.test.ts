import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IPropertyService } from '$lib/services/property/property-service';
import { PropertyListModule } from './property-list-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

describe('PropertyListModule', () => {
  it('returns properties from service', async () => {
    const properties = [
      { id: 'pr1', address: '1 Main St' },
      { id: 'pr2', address: '2 Oak Ave' },
    ];
    const svc = { getAllProperties: vi.fn().mockResolvedValue(properties) } as unknown as IPropertyService;

    const module = new PropertyListModule(svc);
    await expect(module.propertyListAsync).resolves.toEqual(properties);
    expect(svc.getAllProperties).toHaveBeenCalledOnce();
  });

  it('returns empty array when service throws', async () => {
    const svc = {
      getAllProperties: vi.fn().mockRejectedValue(new Error('network error')),
    } as unknown as IPropertyService;

    const module = new PropertyListModule(svc);
    await expect(module.propertyListAsync).resolves.toEqual([]);
  });
});
