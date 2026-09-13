import { describe, it, expect, vi } from 'vitest';
import { CobrandService } from './cobrand-service';
import type { ICobrandRepository } from '$lib/repositories/cobrand/cobrand-repository';
import type { ICobrandAdminRepository } from '$lib/repositories/cobrand/cobrand-admin-repository';
import type { ICobrandPropertyManagerRepository } from '$lib/repositories/cobrand/cobrand-property-manager-repository';

const mockCobrandRepository: ICobrandRepository = {
	getAllCobrands: vi.fn(),
	getCobrandById: vi.fn(),
	createCobrand: vi.fn(),
	updateCobrand: vi.fn(),
	deleteCobrand: vi.fn()
};

const mockAdminRepository: ICobrandAdminRepository = {
	getCobrandAdminsByCobrandId: vi.fn(),
	getCobrandAdminByUserId: vi.fn(),
	createCobrandAdmin: vi.fn(),
	deleteCobrandAdmin: vi.fn()
};

const mockManagerRepository: ICobrandPropertyManagerRepository = {
	getCobrandPropertyManagersByCobrandId: vi.fn(),
	getCobrandPropertyManagersByPropertyId: vi.fn(),
	createCobrandPropertyManager: vi.fn(),
	deleteCobrandPropertyManager: vi.fn()
};

describe('CobrandService.validateCobrandName', () => {
	const cobrandService = new CobrandService(
		mockCobrandRepository,
		mockAdminRepository,
		mockManagerRepository
	);

	it('passes for a valid name', () => {
		expect(() => cobrandService.validateCobrandName('Acme Real Estate')).not.toThrow();
	});

	it('throws for empty string', () => {
		expect(() => cobrandService.validateCobrandName('')).toThrow('cobrand name is required');
	});

	it('throws for whitespace-only string', () => {
		expect(() => cobrandService.validateCobrandName('   ')).toThrow('cobrand name is required');
	});
});

describe('CobrandService.validateAdminUniqueness', () => {
	const cobrandService = new CobrandService(
		mockCobrandRepository,
		mockAdminRepository,
		mockManagerRepository
	);

	it('passes when user is not already an admin', () => {
		const admins = [
			{ id: 'a1', user: 'u1', cobrand: 'c1', approved: true },
			{ id: 'a2', user: 'u2', cobrand: 'c1', approved: true }
		];
		expect(() => cobrandService.validateAdminUniqueness(admins, 'u3')).not.toThrow();
	});

	it('passes for empty admin list', () => {
		expect(() => cobrandService.validateAdminUniqueness([], 'u1')).not.toThrow();
	});

	it('throws when user is already an admin', () => {
		const admins = [{ id: 'a1', user: 'u1', cobrand: 'c1', approved: true }];
		expect(() => cobrandService.validateAdminUniqueness(admins, 'u1')).toThrow(
			'user is already an admin of this cobrand'
		);
	});
});

describe('CobrandService.getAdminRecordForUserId', () => {
	it('returns the admin record from the repository when found', async () => {
		const admin = { id: 'a1', user: 'u1', cobrand: 'c1', approved: true };
		const adminRepositoryFindingTheUser = {
			...mockAdminRepository,
			getCobrandAdminByUserId: vi.fn().mockResolvedValue(admin)
		};
		const cobrandService = new CobrandService(
			mockCobrandRepository,
			adminRepositoryFindingTheUser,
			mockManagerRepository
		);

		await expect(cobrandService.getAdminRecordForUserId('u1')).resolves.toEqual(admin);
		expect(adminRepositoryFindingTheUser.getCobrandAdminByUserId).toHaveBeenCalledWith('u1');
	});

	it('returns null when the user has no admin record', async () => {
		const adminRepositoryFindingNobody = {
			...mockAdminRepository,
			getCobrandAdminByUserId: vi.fn().mockResolvedValue(null)
		};
		const cobrandService = new CobrandService(
			mockCobrandRepository,
			adminRepositoryFindingNobody,
			mockManagerRepository
		);

		await expect(cobrandService.getAdminRecordForUserId('u2')).resolves.toBeNull();
	});
});
