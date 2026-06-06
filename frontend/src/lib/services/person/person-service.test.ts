import { describe, it, expect, vi } from 'vitest';
import type { IPersonRepository } from '$lib/repositories/person/person-repository';
import type { IPersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
import type { ITenantRepository } from '$lib/repositories/property/tenant-repository';
import type { IHouseholdRepository } from '$lib/repositories/property/household-repository';
import type { IAgentRepository } from '$lib/repositories/agent/agent-repository';
import { PersonService } from './person-service';

const makeRepos = (overrides: {
  personRepo?: Partial<IPersonRepository>;
  ppoRepo?: Partial<IPersonPropertyOwnerRepository>;
  tenantRepo?: Partial<ITenantRepository>;
  householdRepo?: Partial<IHouseholdRepository>;
  agentRepo?: Partial<IAgentRepository>;
} = {}) => ({
  personRepo: {
    getAll: vi.fn().mockResolvedValue([]),
    getById: vi.fn().mockResolvedValue(null),
    getByUserId: vi.fn().mockResolvedValue(null),
    create: vi.fn().mockResolvedValue({ id: 'p1', name: '', DOB: '', user: '', profileImage: '' }),
    update: vi.fn().mockResolvedValue({ id: 'p1', name: '', DOB: '', user: '', profileImage: '' }),
    delete: vi.fn().mockResolvedValue(undefined),
    ...overrides.personRepo,
  } as unknown as IPersonRepository,
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
  agentRepo: {
    getAll: vi.fn().mockResolvedValue([]),
    getById: vi.fn().mockResolvedValue(null),
    getByPersonId: vi.fn().mockResolvedValue([]),
    create: vi.fn().mockResolvedValue({ id: 'ag1', person: 'p1', cobrand: 'c1' }),
    delete: vi.fn().mockResolvedValue(undefined),
    ...overrides.agentRepo,
  } as unknown as IAgentRepository,
});

const makeSvc = (overrides = {}) => {
  const repos = makeRepos(overrides);
  return new PersonService(
    repos.personRepo,
    repos.ppoRepo,
    repos.tenantRepo,
    repos.householdRepo,
    repos.agentRepo,
  );
};

describe('PersonService.validatePersonName', () => {
  it('passes for a valid name', () => {
    expect(() => makeSvc().validatePersonName('Alice')).not.toThrow();
  });

  it('throws for empty string', () => {
    expect(() => makeSvc().validatePersonName('')).toThrow('person name is required');
  });

  it('throws for whitespace-only string', () => {
    expect(() => makeSvc().validatePersonName('   ')).toThrow('person name is required');
  });
});

describe('PersonService.validatePersonDOB', () => {
  it('passes for a valid date', () => {
    expect(() => makeSvc().validatePersonDOB('1990-01-01')).not.toThrow();
  });

  it('throws for empty string', () => {
    expect(() => makeSvc().validatePersonDOB('')).toThrow('date of birth is required');
  });
});

describe('PersonService.getRolesForPerson', () => {
  it('returns empty array when person has no roles', async () => {
    const svc = makeSvc();
    await expect(svc.getRolesForPerson('p1')).resolves.toEqual([]);
  });

  it('returns Owner when person has personPropertyOwner records', async () => {
    const svc = makeSvc({
      ppoRepo: {
        getByPersonId: vi.fn().mockResolvedValue([{ id: 'ppo1', person: 'p1', propertyOwner: 'po1' }]),
      },
    });
    const roles = await svc.getRolesForPerson('p1');
    expect(roles).toContain('Owner');
  });

  it('returns Tenant when person has tenant records', async () => {
    const svc = makeSvc({
      tenantRepo: {
        getByPersonId: vi.fn().mockResolvedValue([{ id: 't1', person: 'p1', property: 'pr1' }]),
      },
    });
    const roles = await svc.getRolesForPerson('p1');
    expect(roles).toContain('Tenant');
  });

  it('returns Household when person has household records', async () => {
    const svc = makeSvc({
      householdRepo: {
        getByPersonId: vi.fn().mockResolvedValue([{ id: 'h1', person: 'p1', property: 'pr1' }]),
      },
    });
    const roles = await svc.getRolesForPerson('p1');
    expect(roles).toContain('Household');
  });

  it('returns Agent when person has agent records', async () => {
    const svc = makeSvc({
      agentRepo: {
        getByPersonId: vi.fn().mockResolvedValue([{ id: 'ag1', person: 'p1', cobrand: 'c1' }]),
      },
    });
    const roles = await svc.getRolesForPerson('p1');
    expect(roles).toContain('Agent');
  });

  it('returns multiple roles when person has multiple roles', async () => {
    const svc = makeSvc({
      tenantRepo: {
        getByPersonId: vi.fn().mockResolvedValue([{ id: 't1', person: 'p1', property: 'pr1' }]),
      },
      householdRepo: {
        getByPersonId: vi.fn().mockResolvedValue([{ id: 'h1', person: 'p1', property: 'pr2' }]),
      },
    });
    const roles = await svc.getRolesForPerson('p1');
    expect(roles).toContain('Tenant');
    expect(roles).toContain('Household');
  });
});
