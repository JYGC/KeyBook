import { describe, it, expect, vi, beforeEach, beforeAll } from 'vitest';

const getPersonByUserId = vi.fn();
const getAdminRecordForUserId = vi.fn();

vi.mock('$lib/api/backend-client', () => ({
	getBackendClient: () => ({ authStore: { record: { id: 'user-1' } } })
}));

vi.mock('$lib/services/person/person-service', () => ({
	PersonService: vi.fn().mockImplementation(() => ({ getPersonByUserId }))
}));

vi.mock('$lib/services/cobrand/cobrand-service', () => ({
	CobrandService: vi.fn().mockImplementation(() => ({ getAdminRecordForUserId }))
}));

let load: typeof import('./+page').load;

beforeAll(async () => {
	({ load } = await import('./+page'));
});

describe('/user/+page.ts load', () => {
	beforeEach(() => {
		getPersonByUserId.mockReset();
		getAdminRecordForUserId.mockReset();
	});

	it('redirects to the property list when a person is linked', async () => {
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
		expect(getAdminRecordForUserId).not.toHaveBeenCalled();
	});

	it('redirects to cobrands when no person but a linked cobrandAdmins record exists', async () => {
		getPersonByUserId.mockResolvedValue(null);
		getAdminRecordForUserId.mockResolvedValue({
			id: 'a1',
			user: 'user-1',
			cobrand: 'c1',
			approved: false
		});

		await expect(load({} as never)).rejects.toMatchObject({
			status: 303,
			location: '/user/cobrands/'
		});
	});

	it('redirects to account setup when neither a person nor a cobrandAdmins record exists', async () => {
		getPersonByUserId.mockResolvedValue(null);
		getAdminRecordForUserId.mockResolvedValue(null);

		await expect(load({} as never)).rejects.toMatchObject({
			status: 303,
			location: '/user/setup'
		});
	});

	it('propagates the error without redirecting when the person lookup fails', async () => {
		getPersonByUserId.mockRejectedValue(new Error('network error'));

		await expect(load({} as never)).rejects.toThrow('network error');
		expect(getAdminRecordForUserId).not.toHaveBeenCalled();
	});

	it('propagates the error without redirecting when the cobrand admin lookup fails', async () => {
		getPersonByUserId.mockResolvedValue(null);
		getAdminRecordForUserId.mockRejectedValue(new Error('network error'));

		await expect(load({} as never)).rejects.toThrow('network error');
	});
});
