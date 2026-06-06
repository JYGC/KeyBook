import { describe, it, expect, vi } from 'vitest';
import type { IPropertyRepository } from '$lib/repositories/property/property-repository';
import type { IPropertyOwnerRepository } from '$lib/repositories/property/property-owner-repository';
import type { IPersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
import type { ITenantRepository } from '$lib/repositories/property/tenant-repository';
import type { IHouseholdRepository } from '$lib/repositories/property/household-repository';
import type { IPropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';
import { PropertyService } from './property-service';

const makeRepos = (overrides: {
  propertyRepo?: Partial<IPropertyRepository>;
  propertyOwnerRepo?: Partial<IPropertyOwnerRepository>;
  ppoRepo?: Partial<IPersonPropertyOwnerRepository>;
  tenantRepo?: Partial<ITenantRepository>;
  householdRepo?: Partial<IHouseholdRepository>;
  propertyAgentRepo?: Partial<IPropertyAgentRepository>;
} = {}) => ({
  propertyRepo: {
    getAll: vi.fn().mockResolvedValue([]),
    getById: vi.fn().mockResolvedValue({ id: 'pr1', address: '1 Test St' }),
    create: vi.fn().mockResolvedValue({ id: 'pr1', address: '1 Test St' }),
    update: vi.fn().mockResolvedValue({ id: 'pr1', address: '1 Test St' }),
    delete: vi.fn().mockResolvedValue(undefined),
    ...overrides.propertyRepo,
  } as unknown as IPropertyRepository,
  propertyOwnerRepo: {
    getByPropertyId: vi.fn().mockResolvedValue([]),
    create: vi.fn().mockResolvedValue({ id: 'po1', property: 'pr1' }),
    delete: vi.fn().mockResolvedValue(undefined),
    ...overrides.propertyOwnerRepo,
  } as unknown as IPropertyOwnerRepository,
  ppoRepo: {
    getByPersonId: vi.fn().mockResolvedValue([]),
    getByPropertyOwnerId: vi.fn().mockResolvedValue([]),
    create: vi.fn().mockResolvedValue({ id: 'ppo1', person: 'p1', propertyOwner: 'po1' }),
    delete: vi.fn().mockResolvedValue(undefined),
    ...overrides.ppoRepo,
  } as unknown as IPersonPropertyOwnerRepository,
  tenantRepo: {
    getByPersonId: vi.fn().mockResolvedValue([]),
    getByPropertyId: vi.fn().mockResolvedValue([]),
    create: vi.fn().mockResolvedValue({ id: 't1', person: 'p1', property: 'pr1' }),
    delete: vi.fn().mockResolvedValue(undefined),
    ...overrides.tenantRepo,
  } as unknown as ITenantRepository,
  householdRepo: {
    getByPersonId: vi.fn().mockResolvedValue([]),
    getByPropertyId: vi.fn().mockResolvedValue([]),
    create: vi.fn().mockResolvedValue({ id: 'h1', person: 'p1', property: 'pr1' }),
    delete: vi.fn().mockResolvedValue(undefined),
    ...overrides.householdRepo,
  } as unknown as IHouseholdRepository,
  propertyAgentRepo: {
    getByAgentId: vi.fn().mockResolvedValue([]),
    getByPropertyId: vi.fn().mockResolvedValue([]),
    create: vi.fn().mockResolvedValue({ id: 'pa1', agent: 'ag1', property: 'pr1' }),
    delete: vi.fn().mockResolvedValue(undefined),
    ...overrides.propertyAgentRepo,
  } as unknown as IPropertyAgentRepository,
});

const makeSvc = (overrides = {}) => {
  const repos = makeRepos(overrides);
  return new PropertyService(
    repos.propertyRepo,
    repos.propertyOwnerRepo,
    repos.ppoRepo,
    repos.tenantRepo,
    repos.householdRepo,
    repos.propertyAgentRepo,
  );
};

describe('PropertyService.validatePropertyAddress', () => {
  it('passes for a valid address', () => {
    expect(() => makeSvc().validatePropertyAddress('1 Main St')).not.toThrow();
  });

  it('throws for empty string', () => {
    expect(() => makeSvc().validatePropertyAddress('')).toThrow('property address is required');
  });

  it('throws for whitespace-only string', () => {
    expect(() => makeSvc().validatePropertyAddress('   ')).toThrow('property address is required');
  });
});

describe('PropertyService.createPropertyWithOwner', () => {
  it('creates property, propertyOwner, and personPropertyOwner in sequence', async () => {
    const repos = makeRepos({
      propertyRepo: {
        create: vi.fn().mockResolvedValue({ id: 'pr1', address: '1 Test St' }),
      },
      propertyOwnerRepo: {
        create: vi.fn().mockResolvedValue({ id: 'po1', property: 'pr1' }),
      },
      ppoRepo: {
        create: vi.fn().mockResolvedValue({ id: 'ppo1', person: 'p1', propertyOwner: 'po1' }),
      },
    });
    const svc = new PropertyService(
      repos.propertyRepo,
      repos.propertyOwnerRepo,
      repos.ppoRepo,
      repos.tenantRepo,
      repos.householdRepo,
      repos.propertyAgentRepo,
    );

    const result = await svc.createPropertyWithOwner('1 Test St', 'p1');
    expect(repos.propertyRepo.create).toHaveBeenCalledWith('1 Test St');
    expect(repos.propertyOwnerRepo.create).toHaveBeenCalledWith('pr1');
    expect(repos.ppoRepo.create).toHaveBeenCalledWith('p1', 'po1');
    expect(result.id).toBe('pr1');
  });
});
