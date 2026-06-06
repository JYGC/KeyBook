import PocketBase from 'pocketbase';
import type { ICobrandAdminModel } from '$lib/models/cobrand-models';

export interface ICobrandAdminRepository {
  getByCobrandId(cobrandId: string): Promise<ICobrandAdminModel[]>;
  create(userId: string, cobrandId: string): Promise<ICobrandAdminModel>;
  delete(id: string): Promise<void>;
}

export class CobrandAdminRepository implements ICobrandAdminRepository {
  constructor(private readonly pb: PocketBase) {}

  async getByCobrandId(cobrandId: string): Promise<ICobrandAdminModel[]> {
    const all = await this.pb.collection('cobrandAdmins').getFullList<ICobrandAdminModel>();
    return all.filter((a) => a.cobrand === cobrandId);
  }

  async create(userId: string, cobrandId: string): Promise<ICobrandAdminModel> {
    return await this.pb.collection('cobrandAdmins').create<ICobrandAdminModel>({
      user: userId,
      cobrand: cobrandId,
    });
  }

  async delete(id: string): Promise<void> {
    await this.pb.collection('cobrandAdmins').delete(id);
  }
}
