import { describe, it, expect, vi } from 'vitest';
import { ItemService } from './item-service';
import type { IItemRepository } from '$lib/repositories/item/item-repository';
import type { IEntryDeviceRepository } from '$lib/repositories/item/entry-device-repository';

const mockItemRepository: IItemRepository = {
  getAllItems: vi.fn(),
  getItemById: vi.fn(),
  createItem: vi.fn(),
  updateItem: vi.fn(),
  deleteItem: vi.fn(),
};

const mockEntryDeviceRepository: IEntryDeviceRepository = {
  getEntryDeviceById: vi.fn(),
  getEntryDeviceByItemId: vi.fn(),
  createEntryDevice: vi.fn(),
  updateEntryDevice: vi.fn(),
  deleteEntryDevice: vi.fn(),
};

describe('ItemService.validateItemName', () => {
  const itemService = new ItemService(mockItemRepository, mockEntryDeviceRepository);

  it('passes for a valid name', () => {
    expect(() => itemService.validateItemName('Front Door Key')).not.toThrow();
  });

  it('throws for empty string', () => {
    expect(() => itemService.validateItemName('')).toThrow('item name is required');
  });

  it('throws for whitespace-only string', () => {
    expect(() => itemService.validateItemName('   ')).toThrow('item name is required');
  });
});

describe('ItemService.validateEntryDeviceTransition', () => {
  it.each([
    ['active stays active', 'None', 'None', false],
    ['active (empty) stays active', '', 'None', false],
    ['active to defunct', 'None', 'Lost', false],
    ['active to defunct Damaged', 'None', 'Damaged', false],
    ['defunct stays defunct same reason', 'Lost', 'Lost', false],
    ['defunct changes defunct reason', 'Lost', 'Stolen', false],
    ['defunct to active — blocked', 'Lost', 'None', true],
    ['defunct to empty — blocked', 'Damaged', '', true],
  ])('%s: current=%s next=%s throws=%s', (_label, current, next, shouldThrow) => {
    const itemService = new ItemService(mockItemRepository, mockEntryDeviceRepository);
    if (shouldThrow) {
      expect(() => itemService.validateEntryDeviceTransition(current, next)).toThrow(
        'cannot reactivate a defunct entry device'
      );
    } else {
      expect(() => itemService.validateEntryDeviceTransition(current, next)).not.toThrow();
    }
  });
});
