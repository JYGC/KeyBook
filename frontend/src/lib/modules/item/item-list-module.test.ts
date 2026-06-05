import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IItemService } from '$lib/services/item/item-service';
import { ItemListModule } from './item-list-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

describe('ItemListModule', () => {
  it('returns items from service', async () => {
    const items = [
      { id: '1', name: 'Front Door Key', description: 'Key to front door' },
      { id: '2', name: 'Mailbox Key', description: 'Key to mailbox' },
    ];
    const svc = { getAllItems: vi.fn().mockResolvedValue(items) } as unknown as IItemService;

    const module = new ItemListModule(svc);
    await expect(module.itemListAsync).resolves.toEqual(items);
    expect(svc.getAllItems).toHaveBeenCalledOnce();
  });

  it('returns empty array when service throws', async () => {
    const svc = {
      getAllItems: vi.fn().mockRejectedValue(new Error('network error')),
    } as unknown as IItemService;

    const module = new ItemListModule(svc);
    await expect(module.itemListAsync).resolves.toEqual([]);
  });
});
