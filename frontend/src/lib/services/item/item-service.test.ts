import { describe, it, expect, vi } from 'vitest';
import { ItemService } from './item-service';
import type { IItemRepository } from '$lib/repositories/item/item-repository';
import type { IEntryDeviceRepository } from '$lib/repositories/item/entry-device-repository';

const mockItemRepo: IItemRepository = {
  getAll: vi.fn(),
  getById: vi.fn(),
  create: vi.fn(),
  update: vi.fn(),
  delete: vi.fn(),
};

const mockEntryDeviceRepo: IEntryDeviceRepository = {
  getById: vi.fn(),
  getByItemId: vi.fn(),
  create: vi.fn(),
  update: vi.fn(),
  delete: vi.fn(),
};

describe('ItemService.validateItemName', () => {
  const svc = new ItemService(mockItemRepo, mockEntryDeviceRepo);

  it('passes for a valid name', () => {
    expect(() => svc.validateItemName('Front Door Key')).not.toThrow();
  });

  it('throws for empty string', () => {
    expect(() => svc.validateItemName('')).toThrow('item name is required');
  });

  it('throws for whitespace-only string', () => {
    expect(() => svc.validateItemName('   ')).toThrow('item name is required');
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
    const svc = new ItemService(mockItemRepo, mockEntryDeviceRepo);
    if (shouldThrow) {
      expect(() => svc.validateEntryDeviceTransition(current, next)).toThrow(
        'cannot reactivate a defunct entry device'
      );
    } else {
      expect(() => svc.validateEntryDeviceTransition(current, next)).not.toThrow();
    }
  });
});
