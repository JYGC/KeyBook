import type { ICobrandService } from '$lib/services/cobrand/cobrand-service';
import type { ICobrandDetailModule } from '$lib/modules/interfaces';
import type {
  ICobrandModel,
  ICobrandAdminModel,
  ICobrandPropertyManagerModel,
} from '$lib/models/cobrand-models';

export class CobrandDetailModule implements ICobrandDetailModule {
  private readonly __cobrandService: ICobrandService;
  private readonly __cobrandId: string;
  private readonly __backAction: () => void;

  public cobrandAsync: Promise<ICobrandModel | null>;
  public adminsAsync = $state<Promise<ICobrandAdminModel[]>>(Promise.resolve([]));
  public propertyManagersAsync = $state<Promise<ICobrandPropertyManagerModel[]>>(
    Promise.resolve([]),
  );

  get isAdd() {
    return false;
  }

  public callBackAction = () => this.__backAction();

  private refreshAdmins() {
    this.adminsAsync = this.__cobrandService
      .getAdminsForCobrand(this.__cobrandId)
      .catch((ex) => {
        alert(ex);
        return [];
      });
  }

  private refreshPropertyManagers() {
    this.propertyManagersAsync = this.__cobrandService
      .getPropertyManagersForCobrand(this.__cobrandId)
      .catch((ex) => {
        alert(ex);
        return [];
      });
  }

  public getSaveCobrandAction = () => async (cobrand: ICobrandModel) => {
    try {
      await this.__cobrandService.updateCobrand(cobrand.id, cobrand.name);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getDeleteCobrandAction = () => async (id: string) => {
    try {
      await this.__cobrandService.deleteCobrand(id);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getAddAdminAction = () => async (cobrandId: string, userId: string) => {
    try {
      await this.__cobrandService.addAdmin(cobrandId, userId);
      this.refreshAdmins();
    } catch (ex) {
      alert(ex);
    }
  };

  public getRemoveAdminAction = () => async (id: string) => {
    try {
      await this.__cobrandService.removeAdmin(id);
      this.refreshAdmins();
    } catch (ex) {
      alert(ex);
    }
  };

  public getAddPropertyManagerAction = () => async (cobrandId: string, propertyId: string) => {
    try {
      await this.__cobrandService.addPropertyManager(cobrandId, propertyId);
      this.refreshPropertyManagers();
    } catch (ex) {
      alert(ex);
    }
  };

  public getRemovePropertyManagerAction = () => async (id: string) => {
    try {
      await this.__cobrandService.removePropertyManager(id);
      this.refreshPropertyManagers();
    } catch (ex) {
      alert(ex);
    }
  };

  constructor(cobrandService: ICobrandService, cobrandId: string, backAction: () => void) {
    this.__cobrandService = cobrandService;
    this.__cobrandId = cobrandId;
    this.__backAction = backAction;

    this.cobrandAsync = this.__cobrandService
      .getCobrandById(cobrandId)
      .catch((ex) => {
        alert(ex);
        return null;
      });

    this.refreshAdmins();
    this.refreshPropertyManagers();
  }
}
