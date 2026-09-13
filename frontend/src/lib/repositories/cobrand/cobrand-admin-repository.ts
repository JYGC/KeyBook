import PocketBase from 'pocketbase';
import type { ICobrandAdminModel } from '$lib/models/cobrand-models';

export interface ICobrandAdminRepository {
	getCobrandAdminsByCobrandId(cobrandId: string): Promise<ICobrandAdminModel[]>;
	getCobrandAdminByUserId(userId: string): Promise<ICobrandAdminModel | null>;
	createCobrandAdmin(userId: string, cobrandId: string): Promise<ICobrandAdminModel>;
	deleteCobrandAdmin(id: string): Promise<void>;
}

export class CobrandAdminRepository implements ICobrandAdminRepository {
	constructor(private readonly backendClient: PocketBase) {}

	async getCobrandAdminsByCobrandId(cobrandId: string): Promise<ICobrandAdminModel[]> {
		const allCobrandAdmins = await this.backendClient.collection('cobrandAdmins').getFullList<ICobrandAdminModel>();
		return allCobrandAdmins.filter((cobrandAdmin) => cobrandAdmin.cobrand === cobrandId);
	}

	async getCobrandAdminByUserId(userId: string): Promise<ICobrandAdminModel | null> {
		const allCobrandAdmins = await this.backendClient.collection('cobrandAdmins').getFullList<ICobrandAdminModel>();
		return allCobrandAdmins.find((cobrandAdmin) => cobrandAdmin.user === userId) ?? null;
	}

	async createCobrandAdmin(userId: string, cobrandId: string): Promise<ICobrandAdminModel> {
		return await this.backendClient.collection('cobrandAdmins').create<ICobrandAdminModel>({
			user: userId,
			cobrand: cobrandId
		});
	}

	async deleteCobrandAdmin(id: string): Promise<void> {
		await this.backendClient.collection('cobrandAdmins').delete(id);
	}
}
