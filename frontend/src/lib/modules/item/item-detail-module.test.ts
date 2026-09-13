import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IItemService } from '$lib/services/item/item-service';
import { ItemDetailModule } from './item-detail-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

const mockItem = { id: 'item1', name: 'Test Key', description: 'A key' };
const mockEntryDevice = { id: 'ed1', item: 'item1', deviceType: 'Key', identifier: 'K-001', defunctReason: 'None' };

describe('ItemDetailModule', () => {
  it('isAdd returns false', () => {
    const itemService = { getItemWithEntryDevice: vi.fn().mockResolvedValue({ item: mockItem, entryDevice: null }) } as unknown as IItemService;
    const module = new ItemDetailModule(itemService, 'item1', vi.fn());
    expect(module.isAdd).toBe(false);
  });

  it('callBackAction calls the provided callback', () => {
    const itemService = { getItemWithEntryDevice: vi.fn().mockResolvedValue({ item: mockItem, entryDevice: null }) } as unknown as IItemService;
    const backAction = vi.fn();
    const module = new ItemDetailModule(itemService, 'item1', backAction);
    module.callBackAction();
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('resolves item from service', async () => {
    const itemService = {
      getItemWithEntryDevice: vi.fn().mockResolvedValue({ item: mockItem, entryDevice: null }),
    } as unknown as IItemService;
    const module = new ItemDetailModule(itemService, 'item1', vi.fn());
    await expect(module.itemAsync).resolves.toEqual(mockItem);
  });

  it('resolves entryDevice when one exists', async () => {
    const itemService = {
      getItemWithEntryDevice: vi.fn().mockResolvedValue({ item: mockItem, entryDevice: mockEntryDevice }),
    } as unknown as IItemService;
    const module = new ItemDetailModule(itemService, 'item1', vi.fn());
    await expect(module.entryDeviceAsync).resolves.toEqual(mockEntryDevice);
  });

  it('resolves null entryDevice when none exists', async () => {
    const itemService = {
      getItemWithEntryDevice: vi.fn().mockResolvedValue({ item: mockItem, entryDevice: null }),
    } as unknown as IItemService;
    const module = new ItemDetailModule(itemService, 'item1', vi.fn());
    await expect(module.entryDeviceAsync).resolves.toBeNull();
  });

  it('getSaveItemAction returns a callable function', () => {
    const itemService = { getItemWithEntryDevice: vi.fn().mockResolvedValue({ item: mockItem, entryDevice: null }) } as unknown as IItemService;
    const module = new ItemDetailModule(itemService, 'item1', vi.fn());
    expect(typeof module.getSaveItemAction()).toBe('function');
  });

  it('getDeleteItemAction returns a callable function', () => {
    const itemService = { getItemWithEntryDevice: vi.fn().mockResolvedValue({ item: mockItem, entryDevice: null }) } as unknown as IItemService;
    const module = new ItemDetailModule(itemService, 'item1', vi.fn());
    expect(typeof module.getDeleteItemAction()).toBe('function');
  });

  it('getSaveItemAction calls service and navigates back', async () => {
    const backAction = vi.fn();
    const itemService = {
      getItemWithEntryDevice: vi.fn().mockResolvedValue({ item: mockItem, entryDevice: null }),
      updateItem: vi.fn().mockResolvedValue(undefined),
    } as unknown as IItemService;
    const module = new ItemDetailModule(itemService, 'item1', backAction);
    await module.getSaveItemAction()(mockItem);
    expect(itemService.updateItem).toHaveBeenCalledWith('item1', 'Test Key', 'A key');
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('getDeleteItemAction calls service and navigates back', async () => {
    const backAction = vi.fn();
    const itemService = {
      getItemWithEntryDevice: vi.fn().mockResolvedValue({ item: mockItem, entryDevice: null }),
      deleteItem: vi.fn().mockResolvedValue(undefined),
    } as unknown as IItemService;
    const module = new ItemDetailModule(itemService, 'item1', backAction);
    await module.getDeleteItemAction()!('item1');
    expect(itemService.deleteItem).toHaveBeenCalledWith('item1');
    expect(backAction).toHaveBeenCalledOnce();
  });
});
