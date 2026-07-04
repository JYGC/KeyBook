import { describe, it, expect } from 'vitest';
import { vi } from 'vitest';
import type PocketBase from 'pocketbase';
import { RegisterModule } from './register-module.svelte';

const makeBackendClient = (create: ReturnType<typeof vi.fn>) =>
	({
		collection: vi.fn().mockReturnValue({ create })
	}) as unknown as PocketBase;

describe('RegisterModule', () => {
	it('starts with empty email, password, passwordConfirm, and error', () => {
		const module = new RegisterModule(makeBackendClient(vi.fn()));
		expect(module.email).toBe('');
		expect(module.password).toBe('');
		expect(module.passwordConfirm).toBe('');
		expect(module.error).toBe('');
	});

	it('callApi creates the user without a name field and returns true on success', async () => {
		const create = vi.fn().mockResolvedValue({ id: 'u1' });
		const backendClient = makeBackendClient(create);
		const module = new RegisterModule(backendClient);
		module.email = 'alice@example.com';
		module.password = 'Password123!';
		module.passwordConfirm = 'Password123!';

		const result = await module.callApi();

		expect(result).toBe(true);
		expect(backendClient.collection).toHaveBeenCalledWith('users');
		expect(create).toHaveBeenCalledWith({
			password: 'Password123!',
			passwordConfirm: 'Password123!',
			email: 'alice@example.com',
			emailVisibility: false
		});
	});

	it('callApi surfaces the error and returns false on failure', async () => {
		// Shaped like a PocketBase ClientResponseError: `message` on the underlying
		// Error is a generic constant, so the real detail lives in `response`, which
		// is what JSON.stringify(ex) actually surfaces to the user.
		const failure = {
			status: 400,
			response: { email: { message: 'email already taken' } }
		};
		const create = vi.fn().mockRejectedValue(failure);
		const module = new RegisterModule(makeBackendClient(create));
		module.email = 'taken@example.com';
		module.password = 'Password123!';
		module.passwordConfirm = 'Password123!';

		const result = await module.callApi();

		expect(result).toBe(false);
		expect(module.error).toContain('email already taken');
	});
});
