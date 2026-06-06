import type { ICobrandService } from '$lib/services/cobrand/cobrand-service';
import type { ICobrandListModule } from '$lib/modules/interfaces';
import type { ICobrandModel } from '$lib/models/cobrand-models';

export class CobrandListModule implements ICobrandListModule {
  private readonly __cobrandService: ICobrandService;

  public cobrandListAsync = $derived.by<Promise<ICobrandModel[]>>(async () => {
    try {
      return await this.__cobrandService.getAllCobrands();
    } catch (ex) {
      alert(ex);
      return [];
    }
  });

  constructor(cobrandService: ICobrandService) {
    this.__cobrandService = cobrandService;
  }
}
