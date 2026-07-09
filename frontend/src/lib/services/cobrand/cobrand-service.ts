import type { ICobrandRepository } from '$lib/repositories/cobrand/cobrand-repository';
import type { ICobrandAdminRepository } from '$lib/repositories/cobrand/cobrand-admin-repository';
import type { ICobrandPropertyManagerRepository } from '$lib/repositories/cobrand/cobrand-property-manager-repository';
import type {
	ICobrandModel,
	ICobrandAdminModel,
	ICobrandPropertyManagerModel
} from '$lib/models/cobrand-models';

export interface ICobrandService {
	validateCobrandName(name: string): void;
	validateAdminUniqueness(existingAdmins: ICobrandAdminModel[], userId: string): void;
	getAllCobrands(): Promise<ICobrandModel[]>;
	getCobrandById(id: string): Promise<ICobrandModel>;
	createCobrand(name: string): Promise<ICobrandModel>;
	updateCobrand(id: string, name: string): Promise<void>;
	deleteCobrand(id: string): Promise<void>;
	getAdminsForCobrand(cobrandId: string): Promise<ICobrandAdminModel[]>;
	getAdminRecordForUserId(userId: string): Promise<ICobrandAdminModel | null>;
	addAdmin(cobrandId: string, userId: string): Promise<ICobrandAdminModel>;
	removeAdmin(id: string): Promise<void>;
	getPropertyManagersForCobrand(cobrandId: string): Promise<ICobrandPropertyManagerModel[]>;
	addPropertyManager(cobrandId: string, propertyId: string): Promise<ICobrandPropertyManagerModel>;
	removePropertyManager(id: string): Promise<void>;
}

export class CobrandService implements ICobrandService {
	constructor(
		private readonly cobrandRepo: ICobrandRepository,
		private readonly cobrandAdminRepo: ICobrandAdminRepository,
		private readonly cobrandPropertyManagerRepo: ICobrandPropertyManagerRepository
	) {}

	validateCobrandName(name: string): void {
		if (!name || name.trim() === '') {
			throw new Error('cobrand name is required');
		}
	}

	validateAdminUniqueness(existingAdmins: ICobrandAdminModel[], userId: string): void {
		if (existingAdmins.some((a) => a.user === userId)) {
			throw new Error('user is already an admin of this cobrand');
		}
	}

	async getAllCobrands(): Promise<ICobrandModel[]> {
		return await this.cobrandRepo.getAll();
	}

	async getCobrandById(id: string): Promise<ICobrandModel> {
		return await this.cobrandRepo.getById(id);
	}

	async createCobrand(name: string): Promise<ICobrandModel> {
		this.validateCobrandName(name);
		return await this.cobrandRepo.create(name);
	}

	async updateCobrand(id: string, name: string): Promise<void> {
		this.validateCobrandName(name);
		await this.cobrandRepo.update(id, name);
	}

	async deleteCobrand(id: string): Promise<void> {
		await this.cobrandRepo.delete(id);
	}

	async getAdminsForCobrand(cobrandId: string): Promise<ICobrandAdminModel[]> {
		return await this.cobrandAdminRepo.getByCobrandId(cobrandId);
	}

	async getAdminRecordForUserId(userId: string): Promise<ICobrandAdminModel | null> {
		return await this.cobrandAdminRepo.getByUserId(userId);
	}

	async addAdmin(cobrandId: string, userId: string): Promise<ICobrandAdminModel> {
		const existing = await this.cobrandAdminRepo.getByCobrandId(cobrandId);
		this.validateAdminUniqueness(existing, userId);
		return await this.cobrandAdminRepo.create(userId, cobrandId);
	}

	async removeAdmin(id: string): Promise<void> {
		await this.cobrandAdminRepo.delete(id);
	}

	async getPropertyManagersForCobrand(cobrandId: string): Promise<ICobrandPropertyManagerModel[]> {
		return await this.cobrandPropertyManagerRepo.getByCobrandId(cobrandId);
	}

	async addPropertyManager(
		cobrandId: string,
		propertyId: string
	): Promise<ICobrandPropertyManagerModel> {
		return await this.cobrandPropertyManagerRepo.create(cobrandId, propertyId);
	}

	async removePropertyManager(id: string): Promise<void> {
		await this.cobrandPropertyManagerRepo.delete(id);
	}
}
