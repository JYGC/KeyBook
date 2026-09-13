import { describe, it, expect, vi } from 'vitest';
import type { IPersonRepository } from '$lib/repositories/person/person-repository';
import type { IPersonPropertyOwnerRepository } from '$lib/repositories/person/person-property-owner-repository';
import type { ITenantRepository } from '$lib/repositories/property/tenant-repository';
import type { IHouseholdRepository } from '$lib/repositories/property/household-repository';
import type { IAgentRepository } from '$lib/repositories/agent/agent-repository';
import { PersonService } from './person-service';

const makeRepositories = (
	overrides: {
		personRepository?: Partial<IPersonRepository>;
		personPropertyOwnerRepository?: Partial<IPersonPropertyOwnerRepository>;
		tenantRepository?: Partial<ITenantRepository>;
		householdRepository?: Partial<IHouseholdRepository>;
		agentRepository?: Partial<IAgentRepository>;
	} = {}
) => ({
	personRepository: {
		getAllPersons: vi.fn().mockResolvedValue([]),
		getPersonById: vi.fn().mockResolvedValue(null),
		getPersonByUserId: vi.fn().mockResolvedValue(null),
		createPerson: vi.fn().mockResolvedValue({ id: 'p1', name: '', DOB: '', user: '', profileImage: '' }),
		updatePerson: vi.fn().mockResolvedValue({ id: 'p1', name: '', DOB: '', user: '', profileImage: '' }),
		deletePerson: vi.fn().mockResolvedValue(undefined),
		...overrides.personRepository
	} as unknown as IPersonRepository,
	personPropertyOwnerRepository: {
		getPersonPropertyOwnersByPersonId: vi.fn().mockResolvedValue([]),
		getPersonPropertyOwnersByPropertyOwnerId: vi.fn().mockResolvedValue([]),
		createPersonPropertyOwner: vi.fn().mockResolvedValue({ id: 'ppo1', person: 'p1', propertyOwner: 'po1' }),
		deletePersonPropertyOwner: vi.fn().mockResolvedValue(undefined),
		...overrides.personPropertyOwnerRepository
	} as unknown as IPersonPropertyOwnerRepository,
	tenantRepository: {
		getTenantsByPersonId: vi.fn().mockResolvedValue([]),
		getTenantsByPropertyId: vi.fn().mockResolvedValue([]),
		createTenant: vi.fn().mockResolvedValue({ id: 't1', person: 'p1', property: 'pr1' }),
		deleteTenant: vi.fn().mockResolvedValue(undefined),
		...overrides.tenantRepository
	} as unknown as ITenantRepository,
	householdRepository: {
		getHouseholdsByPersonId: vi.fn().mockResolvedValue([]),
		getHouseholdsByPropertyId: vi.fn().mockResolvedValue([]),
		createHousehold: vi.fn().mockResolvedValue({ id: 'h1', person: 'p1', property: 'pr1' }),
		deleteHousehold: vi.fn().mockResolvedValue(undefined),
		...overrides.householdRepository
	} as unknown as IHouseholdRepository,
	agentRepository: {
		getAllAgents: vi.fn().mockResolvedValue([]),
		getAgentById: vi.fn().mockResolvedValue(null),
		getAgentsByPersonId: vi.fn().mockResolvedValue([]),
		createAgent: vi.fn().mockResolvedValue({ id: 'ag1', person: 'p1', cobrand: 'c1' }),
		deleteAgent: vi.fn().mockResolvedValue(undefined),
		...overrides.agentRepository
	} as unknown as IAgentRepository
});

const makePersonService = (overrides = {}) => {
	const repositories = makeRepositories(overrides);
	return new PersonService(
		repositories.personRepository,
		repositories.personPropertyOwnerRepository,
		repositories.tenantRepository,
		repositories.householdRepository,
		repositories.agentRepository
	);
};

describe('PersonService.validatePersonName', () => {
	it('passes for a valid name', () => {
		expect(() => makePersonService().validatePersonName('Alice')).not.toThrow();
	});

	it('throws for empty string', () => {
		expect(() => makePersonService().validatePersonName('')).toThrow('person name is required');
	});

	it('throws for whitespace-only string', () => {
		expect(() => makePersonService().validatePersonName('   ')).toThrow('person name is required');
	});
});

describe('PersonService.createPerson', () => {
	it('creates a person without a userId when none is given', async () => {
		const createPerson = vi
			.fn()
			.mockResolvedValue({
				id: 'p1',
				name: 'Alice',
				DOB: '1990-01-01',
				user: '',
				profileImage: ''
			});
		const personService = makePersonService({ personRepository: { createPerson } });

		await personService.createPerson('Alice', '1990-01-01');

		expect(createPerson).toHaveBeenCalledWith('Alice', '1990-01-01', undefined);
	});

	it('passes the userId through to the repository when linking to the authenticated user', async () => {
		const createPerson = vi
			.fn()
			.mockResolvedValue({
				id: 'p1',
				name: 'Alice',
				DOB: '1990-01-01',
				user: 'u1',
				profileImage: ''
			});
		const personService = makePersonService({ personRepository: { createPerson } });

		const result = await personService.createPerson('Alice', '1990-01-01', 'u1');

		expect(createPerson).toHaveBeenCalledWith('Alice', '1990-01-01', 'u1');
		expect(result.user).toBe('u1');
	});

	it('still validates name and DOB before creating a linked person', async () => {
		const personService = makePersonService();
		await expect(personService.createPerson('', '1990-01-01', 'u1')).rejects.toThrow(
			'person name is required'
		);
	});
});

describe('PersonService.validatePersonDOB', () => {
	it('passes for a valid date', () => {
		expect(() => makePersonService().validatePersonDOB('1990-01-01')).not.toThrow();
	});

	it('throws for empty string', () => {
		expect(() => makePersonService().validatePersonDOB('')).toThrow('date of birth is required');
	});
});

describe('PersonService.getRolesForPerson', () => {
	it('returns empty array when person has no roles', async () => {
		const personService = makePersonService();
		await expect(personService.getRolesForPerson('p1')).resolves.toEqual([]);
	});

	it('returns Owner when person has personPropertyOwner records', async () => {
		const personService = makePersonService({
			personPropertyOwnerRepository: {
				getPersonPropertyOwnersByPersonId: vi
					.fn()
					.mockResolvedValue([{ id: 'ppo1', person: 'p1', propertyOwner: 'po1' }])
			}
		});
		const roles = await personService.getRolesForPerson('p1');
		expect(roles).toContain('Owner');
	});

	it('returns Tenant when person has tenant records', async () => {
		const personService = makePersonService({
			tenantRepository: {
				getTenantsByPersonId: vi.fn().mockResolvedValue([{ id: 't1', person: 'p1', property: 'pr1' }])
			}
		});
		const roles = await personService.getRolesForPerson('p1');
		expect(roles).toContain('Tenant');
	});

	it('returns Household when person has household records', async () => {
		const personService = makePersonService({
			householdRepository: {
				getHouseholdsByPersonId: vi.fn().mockResolvedValue([{ id: 'h1', person: 'p1', property: 'pr1' }])
			}
		});
		const roles = await personService.getRolesForPerson('p1');
		expect(roles).toContain('Household');
	});

	it('returns Agent when person has agent records', async () => {
		const personService = makePersonService({
			agentRepository: {
				getAgentsByPersonId: vi.fn().mockResolvedValue([{ id: 'ag1', person: 'p1', cobrand: 'c1' }])
			}
		});
		const roles = await personService.getRolesForPerson('p1');
		expect(roles).toContain('Agent');
	});

	it('returns multiple roles when person has multiple roles', async () => {
		const personService = makePersonService({
			tenantRepository: {
				getTenantsByPersonId: vi.fn().mockResolvedValue([{ id: 't1', person: 'p1', property: 'pr1' }])
			},
			householdRepository: {
				getHouseholdsByPersonId: vi.fn().mockResolvedValue([{ id: 'h1', person: 'p1', property: 'pr2' }])
			}
		});
		const roles = await personService.getRolesForPerson('p1');
		expect(roles).toContain('Tenant');
		expect(roles).toContain('Household');
	});
});
