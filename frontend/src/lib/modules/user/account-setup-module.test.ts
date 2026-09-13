import { describe, it, expect, vi } from 'vitest';
import type { IPersonService } from '$lib/services/person/person-service';
import type { ICobrandService } from '$lib/services/cobrand/cobrand-service';
import { AccountSetupModule } from './account-setup-module.svelte';

const makePersonSvc = (overrides: Partial<IPersonService> = {}) =>
	({
		getPersonByUserId: vi.fn().mockResolvedValue(null),
		...overrides
	}) as unknown as IPersonService;

const makeCobrandSvc = (overrides: Partial<ICobrandService> = {}) =>
	({
		getAdminRecordForUserId: vi.fn().mockResolvedValue(null),
		...overrides
	}) as unknown as ICobrandService;

describe('AccountSetupModule', () => {
	it('shows the choice when neither a person nor a cobrand admin record is linked', async () => {
		const redirectAction = vi.fn();
		const module = new AccountSetupModule(makePersonSvc(), makeCobrandSvc(), 'u1', redirectAction);

		await expect(module.showChoiceAsync).resolves.toBe(true);
		expect(redirectAction).not.toHaveBeenCalled();
	});

	it('redirects to the property list and hides the choice when a person is linked', async () => {
		const redirectAction = vi.fn();
		const personSvc = makePersonSvc({
			getPersonByUserId: vi
				.fn()
				.mockResolvedValue({
					id: 'p1',
					name: 'Alice',
					DOB: '1990-01-01',
					user: 'u1',
					profileImage: ''
				})
		});
		const cobrandSvc = makeCobrandSvc();
		const module = new AccountSetupModule(personSvc, cobrandSvc, 'u1', redirectAction);

		await expect(module.showChoiceAsync).resolves.toBe(false);
		expect(redirectAction).toHaveBeenCalledWith('/user/properties/list');
		expect(cobrandSvc.getAdminRecordForUserId).not.toHaveBeenCalled();
	});

	it('redirects to cobrands and hides the choice when no person but a cobrand admin record is linked', async () => {
		const redirectAction = vi.fn();
		const cobrandSvc = makeCobrandSvc({
			getAdminRecordForUserId: vi
				.fn()
				.mockResolvedValue({ id: 'a1', user: 'u1', cobrand: 'c1', approved: false })
		});
		const module = new AccountSetupModule(makePersonSvc(), cobrandSvc, 'u1', redirectAction);

		await expect(module.showChoiceAsync).resolves.toBe(false);
		expect(redirectAction).toHaveBeenCalledWith('/user/cobrands/');
	});
});
