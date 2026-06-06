import type { ICobrandService } from '$lib/services/cobrand/cobrand-service';
import type { ICobrandEditorModule } from '$lib/modules/interfaces';
import type { ICobrandModel } from '$lib/models/cobrand-models';

export class CobrandAddModule implements ICobrandEditorModule {
  private readonly __cobrandService: ICobrandService;
  private readonly __currentUserId: string;
  private readonly __backAction: () => void;

  public cobrandAsync = $state<Promise<ICobrandModel | null>>(
    Promise.resolve({ id: '', name: '' }),
  );

  get isAdd() {
    return true;
  }

  public callBackAction = () => this.__backAction();

  public getSaveCobrandAction = () => async (cobrand: ICobrandModel) => {
    try {
      const created = await this.__cobrandService.createCobrand(cobrand.name);
      await this.__cobrandService.addAdmin(created.id, this.__currentUserId);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getDeleteCobrandAction = () => null;

  constructor(cobrandService: ICobrandService, currentUserId: string, backAction: () => void) {
    this.__cobrandService = cobrandService;
    this.__currentUserId = currentUserId;
    this.__backAction = backAction;
  }
}
