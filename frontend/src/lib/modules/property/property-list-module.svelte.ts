import type { IPropertyService } from '$lib/services/property/property-service';
import type { INewPropertyListModule } from '$lib/modules/interfaces';
import type { IPropertyModel } from '$lib/models/property-models';

export class PropertyListModule implements INewPropertyListModule {
  private readonly __propertyService: IPropertyService;

  public propertyListAsync = $derived.by<Promise<IPropertyModel[]>>(async () => {
    try {
      return await this.__propertyService.getAllProperties();
    } catch (ex) {
      alert(ex);
      return [];
    }
  });

  constructor(propertyService: IPropertyService) {
    this.__propertyService = propertyService;
  }
}
