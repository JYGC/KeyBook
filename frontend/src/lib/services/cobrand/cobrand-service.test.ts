import { describe, it, expect, vi } from 'vitest';
import { CobrandService } from './cobrand-service';
import type { ICobrandRepository } from '$lib/repositories/cobrand/cobrand-repository';
import type { ICobrandAdminRepository } from '$lib/repositories/cobrand/cobrand-admin-repository';
import type { ICobrandPropertyManagerRepository } from '$lib/repositories/cobrand/cobrand-property-manager-repository';

const mockCobrandRepo: ICobrandRepository = {
  getAll: vi.fn(),
  getById: vi.fn(),
  create: vi.fn(),
  update: vi.fn(),
  delete: vi.fn(),
};

const mockAdminRepo: ICobrandAdminRepository = {
  getByCobrandId: vi.fn(),
  create: vi.fn(),
  delete: vi.fn(),
};

const mockManagerRepo: ICobrandPropertyManagerRepository = {
  getByCobrandId: vi.fn(),
  getByPropertyId: vi.fn(),
  create: vi.fn(),
  delete: vi.fn(),
};

describe('CobrandService.validateCobrandName', () => {
  const svc = new CobrandService(mockCobrandRepo, mockAdminRepo, mockManagerRepo);

  it('passes for a valid name', () => {
    expect(() => svc.validateCobrandName('Acme Real Estate')).not.toThrow();
  });

  it('throws for empty string', () => {
    expect(() => svc.validateCobrandName('')).toThrow('cobrand name is required');
  });

  it('throws for whitespace-only string', () => {
    expect(() => svc.validateCobrandName('   ')).toThrow('cobrand name is required');
  });
});

describe('CobrandService.validateAdminUniqueness', () => {
  const svc = new CobrandService(mockCobrandRepo, mockAdminRepo, mockManagerRepo);

  it('passes when user is not already an admin', () => {
    const admins = [
      { id: 'a1', user: 'u1', cobrand: 'c1' },
      { id: 'a2', user: 'u2', cobrand: 'c1' },
    ];
    expect(() => svc.validateAdminUniqueness(admins, 'u3')).not.toThrow();
  });

  it('passes for empty admin list', () => {
    expect(() => svc.validateAdminUniqueness([], 'u1')).not.toThrow();
  });

  it('throws when user is already an admin', () => {
    const admins = [
      { id: 'a1', user: 'u1', cobrand: 'c1' },
    ];
    expect(() => svc.validateAdminUniqueness(admins, 'u1')).toThrow(
      'user is already an admin of this cobrand',
    );
  });
});
