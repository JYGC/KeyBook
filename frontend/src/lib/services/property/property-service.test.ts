import { describe, it, expect, vi } from 'vitest';
import type { IPropertyRepository } from '$lib/repositories/property/property-repository';
import type { IPropertyOwnerRepository } from '$lib/repositories/property/property-owner-repository';
import type { IPersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
import type { ITenantRepository } from '$lib/repositories/property/tenant-repository';
import type { IHouseholdRepository } from '$lib/repositories/property/household-repository';
import type { IPropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';
import { PropertyService } from './property-service';

const makeRepositories = (overrides: {
  propertyRepository?: Partial<IPropertyRepository>;
  propertyOwnerRepository?: Partial<IPropertyOwnerRepository>;
  personPropertyOwnerRepository?: Partial<IPersonPropertyOwnerRepository>;
  tenantRepository?: Partial<ITenantRepository>;
  householdRepository?: Partial<IHouseholdRepository>;
  propertyAgentRepository?: Partial<IPropertyAgentRepository>;
} = {}) => ({
  propertyRepository: {
    getAllProperties: vi.fn().mockResolvedValue([]),
    getPropertyById: vi.fn().mockResolvedValue({ id: 'pr1', address: '1 Test St' }),
    createProperty: vi.fn().mockResolvedValue({ id: 'pr1', address: '1 Test St' }),
    updateProperty: vi.fn().mockResolvedValue({ id: 'pr1', address: '1 Test St' }),
    deleteProperty: vi.fn().mockResolvedValue(undefined),
    ...overrides.propertyRepository,
  } as unknown as IPropertyRepository,
  propertyOwnerRepository: {
    getPropertyOwnersByPropertyId: vi.fn().mockResolvedValue([]),
    createPropertyOwner: vi.fn().mockResolvedValue({ id: 'po1', property: 'pr1' }),
    deletePropertyOwner: vi.fn().mockResolvedValue(undefined),
    ...overrides.propertyOwnerRepository,
  } as unknown as IPropertyOwnerRepository,
  personPropertyOwnerRepository: {
    getPersonPropertyOwnersByPersonId: vi.fn().mockResolvedValue([]),
    getPersonPropertyOwnersByPropertyOwnerId: vi.fn().mockResolvedValue([]),
    createPersonPropertyOwner: vi.fn().mockResolvedValue({ id: 'ppo1', person: 'p1', propertyOwner: 'po1' }),
    deletePersonPropertyOwner: vi.fn().mockResolvedValue(undefined),
    ...overrides.personPropertyOwnerRepository,
  } as unknown as IPersonPropertyOwnerRepository,
  tenantRepository: {
    getTenantsByPersonId: vi.fn().mockResolvedValue([]),
    getTenantsByPropertyId: vi.fn().mockResolvedValue([]),
    createTenant: vi.fn().mockResolvedValue({ id: 't1', person: 'p1', property: 'pr1' }),
    deleteTenant: vi.fn().mockResolvedValue(undefined),
    ...overrides.tenantRepository,
  } as unknown as ITenantRepository,
  householdRepository: {
    getHouseholdsByPersonId: vi.fn().mockResolvedValue([]),
    getHouseholdsByPropertyId: vi.fn().mockResolvedValue([]),
    createHousehold: vi.fn().mockResolvedValue({ id: 'h1', person: 'p1', property: 'pr1' }),
    deleteHousehold: vi.fn().mockResolvedValue(undefined),
    ...overrides.householdRepository,
  } as unknown as IHouseholdRepository,
  propertyAgentRepository: {
    getPropertyAgentsByAgentId: vi.fn().mockResolvedValue([]),
    getPropertyAgentsByPropertyId: vi.fn().mockResolvedValue([]),
    createPropertyAgent: vi.fn().mockResolvedValue({ id: 'pa1', agent: 'ag1', property: 'pr1' }),
    deletePropertyAgent: vi.fn().mockResolvedValue(undefined),
    ...overrides.propertyAgentRepository,
  } as unknown as IPropertyAgentRepository,
});

const makePropertyService = (overrides = {}) => {
  const repositories = makeRepositories(overrides);
  return new PropertyService(
    repositories.propertyRepository,
    repositories.propertyOwnerRepository,
    repositories.personPropertyOwnerRepository,
    repositories.tenantRepository,
    repositories.householdRepository,
    repositories.propertyAgentRepository,
  );
};

describe('PropertyService.validatePropertyAddress', () => {
  it('passes for a valid address', () => {
    expect(() => makePropertyService().validatePropertyAddress('1 Main St')).not.toThrow();
  });

  it('throws for empty string', () => {
    expect(() => makePropertyService().validatePropertyAddress('')).toThrow('property address is required');
  });

  it('throws for whitespace-only string', () => {
    expect(() => makePropertyService().validatePropertyAddress('   ')).toThrow('property address is required');
  });
});

describe('PropertyService.createPropertyWithOwner', () => {
  it('creates property, propertyOwner, and personPropertyOwner in sequence', async () => {
    const repositories = makeRepositories({
      propertyRepository: {
        createProperty: vi.fn().mockResolvedValue({ id: 'pr1', address: '1 Test St' }),
      },
      propertyOwnerRepository: {
        createPropertyOwner: vi.fn().mockResolvedValue({ id: 'po1', property: 'pr1' }),
      },
      personPropertyOwnerRepository: {
        createPersonPropertyOwner: vi.fn().mockResolvedValue({ id: 'ppo1', person: 'p1', propertyOwner: 'po1' }),
      },
    });
    const propertyService = new PropertyService(
      repositories.propertyRepository,
      repositories.propertyOwnerRepository,
      repositories.personPropertyOwnerRepository,
      repositories.tenantRepository,
      repositories.householdRepository,
      repositories.propertyAgentRepository,
    );

    const result = await propertyService.createPropertyWithOwner('1 Test St', 'p1');
    expect(repositories.propertyRepository.createProperty).toHaveBeenCalledWith('1 Test St');
    expect(repositories.propertyOwnerRepository.createPropertyOwner).toHaveBeenCalledWith('pr1');
    expect(repositories.personPropertyOwnerRepository.createPersonPropertyOwner).toHaveBeenCalledWith('p1', 'po1');
    expect(result.id).toBe('pr1');
  });
});
