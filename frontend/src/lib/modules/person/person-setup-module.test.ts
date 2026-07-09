import { describe, it, expect, vi } from 'vitest';
import type { IPersonService } from '$lib/services/person/person-service';
import { PersonSetupModule } from './person-setup-module.svelte';

const makeSvc = (overrides: Partial<IPersonService> = {}) =>
	({
		createPerson: vi
			.fn()
			.mockResolvedValue({
				id: 'p1',
				name: 'Alice',
				DOB: '1990-01-01',
				user: 'u1',
				profileImage: ''
			}),
		...overrides
	}) as unknown as IPersonService;

describe('PersonSetupModule', () => {
	it('isAdd returns true', () => {
		const module = new PersonSetupModule(makeSvc(), 'u1', vi.fn());
		expect(module.isAdd).toBe(true);
	});

	it('resolves personAsync to an empty person model', async () => {
		const module = new PersonSetupModule(makeSvc(), 'u1', vi.fn());
		const person = await module.personAsync;
		expect(person).toMatchObject({ id: '', name: '', DOB: '' });
	});

	it('starts with no error', () => {
		const module = new PersonSetupModule(makeSvc(), 'u1', vi.fn());
		expect(module.error).toBe('');
	});

	it('getDeletePersonAction returns null', () => {
		const module = new PersonSetupModule(makeSvc(), 'u1', vi.fn());
		expect(module.getDeletePersonAction()).toBeNull();
	});

	it('getSavePersonAction creates a person linked to the authenticated user and signals success', async () => {
		const back = vi.fn();
		const svc = makeSvc();
		const module = new PersonSetupModule(svc, 'u1', back);

		await module.getSavePersonAction()({
			id: '',
			name: 'Alice',
			DOB: '1990-01-01',
			user: '',
			profileImage: ''
		});

		expect(svc.createPerson).toHaveBeenCalledWith('Alice', '1990-01-01', 'u1');
		expect(back).toHaveBeenCalledOnce();
	});

	it('exposes validation errors as reactive state instead of navigating away', async () => {
		const back = vi.fn();
		const svc = makeSvc({
			createPerson: vi.fn().mockRejectedValue(new Error('person name is required'))
		});
		const module = new PersonSetupModule(svc, 'u1', back);

		await module.getSavePersonAction()({
			id: '',
			name: '',
			DOB: '1990-01-01',
			user: '',
			profileImage: ''
		});

		expect(module.error).toContain('person name is required');
		expect(back).not.toHaveBeenCalled();
	});

	it('exposes creation errors as reactive state', async () => {
		const back = vi.fn();
		const svc = makeSvc({
			createPerson: vi.fn().mockRejectedValue(new Error('server unavailable'))
		});
		const module = new PersonSetupModule(svc, 'u1', back);

		await module.getSavePersonAction()({
			id: '',
			name: 'Alice',
			DOB: '1990-01-01',
			user: '',
			profileImage: ''
		});

		expect(module.error).toContain('server unavailable');
		expect(back).not.toHaveBeenCalled();
	});
});
