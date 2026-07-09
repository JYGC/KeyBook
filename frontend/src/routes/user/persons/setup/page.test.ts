import { describe, it, expect, vi, beforeEach, beforeAll } from 'vitest';

const getPersonByUserId = vi.fn();

vi.mock('$lib/api/backend-client', () => ({
	getBackendClient: () => ({ authStore: { record: { id: 'user-1' } } })
}));

vi.mock('$lib/services/person/person-service', () => ({
	PersonService: vi.fn().mockImplementation(() => ({ getPersonByUserId }))
}));

let load: typeof import('./+page').load;

beforeAll(async () => {
	({ load } = await import('./+page'));
});

describe('/user/persons/setup/+page.ts load', () => {
	beforeEach(() => {
		getPersonByUserId.mockReset();
	});

	it('redirects to the property list when the user already has a linked person', async () => {
		getPersonByUserId.mockResolvedValue({
			id: 'p1',
			name: 'Alice',
			DOB: '1990-01-01',
			user: 'user-1',
			profileImage: ''
		});

		await expect(load({} as never)).rejects.toMatchObject({
			status: 303,
			location: '/user/properties/list'
		});
	});

	it('does not redirect when the user has no linked person', async () => {
		getPersonByUserId.mockResolvedValue(null);

		await expect(load({} as never)).resolves.toBeUndefined();
	});

	it('propagates the error without redirecting when the lookup fails', async () => {
		getPersonByUserId.mockRejectedValue(new Error('network error'));

		await expect(load({} as never)).rejects.toThrow('network error');
	});
});
