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
		private readonly cobrandRepository: ICobrandRepository,
		private readonly cobrandAdminRepository: ICobrandAdminRepository,
		private readonly cobrandPropertyManagerRepository: ICobrandPropertyManagerRepository
	) {}

	validateCobrandName(name: string): void {
		if (!name || name.trim() === '') {
			throw new Error('cobrand name is required');
		}
	}

	validateAdminUniqueness(existingAdmins: ICobrandAdminModel[], userId: string): void {
		if (existingAdmins.some((existingAdmin) => existingAdmin.user === userId)) {
			throw new Error('user is already an admin of this cobrand');
		}
	}

	async getAllCobrands(): Promise<ICobrandModel[]> {
		return await this.cobrandRepository.getAllCobrands();
	}

	async getCobrandById(id: string): Promise<ICobrandModel> {
		return await this.cobrandRepository.getCobrandById(id);
	}

	async createCobrand(name: string): Promise<ICobrandModel> {
		this.validateCobrandName(name);
		return await this.cobrandRepository.createCobrand(name);
	}

	async updateCobrand(id: string, name: string): Promise<void> {
		this.validateCobrandName(name);
		await this.cobrandRepository.updateCobrand(id, name);
	}

	async deleteCobrand(id: string): Promise<void> {
		await this.cobrandRepository.deleteCobrand(id);
	}

	async getAdminsForCobrand(cobrandId: string): Promise<ICobrandAdminModel[]> {
		return await this.cobrandAdminRepository.getCobrandAdminsByCobrandId(cobrandId);
	}

	async getAdminRecordForUserId(userId: string): Promise<ICobrandAdminModel | null> {
		return await this.cobrandAdminRepository.getCobrandAdminByUserId(userId);
	}

	async addAdmin(cobrandId: string, userId: string): Promise<ICobrandAdminModel> {
		const existing = await this.cobrandAdminRepository.getCobrandAdminsByCobrandId(cobrandId);
		this.validateAdminUniqueness(existing, userId);
		return await this.cobrandAdminRepository.createCobrandAdmin(userId, cobrandId);
	}

	async removeAdmin(id: string): Promise<void> {
		await this.cobrandAdminRepository.deleteCobrandAdmin(id);
	}

	async getPropertyManagersForCobrand(cobrandId: string): Promise<ICobrandPropertyManagerModel[]> {
		return await this.cobrandPropertyManagerRepository.getCobrandPropertyManagersByCobrandId(cobrandId);
	}

	async addPropertyManager(
		cobrandId: string,
		propertyId: string
	): Promise<ICobrandPropertyManagerModel> {
		return await this.cobrandPropertyManagerRepository.createCobrandPropertyManager(cobrandId, propertyId);
	}

	async removePropertyManager(id: string): Promise<void> {
		await this.cobrandPropertyManagerRepository.deleteCobrandPropertyManager(id);
	}
}
