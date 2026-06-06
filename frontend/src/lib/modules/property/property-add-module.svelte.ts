import type { IPropertyService } from '$lib/services/property/property-service';
import type { INewPropertyEditorModule } from '$lib/modules/interfaces';
import type { IPropertyModel } from '$lib/models/property-models';

export class PropertyAddModule implements INewPropertyEditorModule {
  private readonly __propertyService: IPropertyService;
  private readonly __currentPersonId: string;
  private readonly __backAction: () => void;

  public propertyAsync = $state<Promise<IPropertyModel | null>>(
    Promise.resolve({ id: '', address: '' }),
  );

  get isAdd() {
    return true;
  }

  public callBackAction = () => this.__backAction();

  public getSavePropertyAction = () => async (property: IPropertyModel) => {
    try {
      await this.__propertyService.createPropertyWithOwner(property.address, this.__currentPersonId);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getDeletePropertyAction = () => null;

  constructor(propertyService: IPropertyService, currentPersonId: string, backAction: () => void) {
    this.__propertyService = propertyService;
    this.__currentPersonId = currentPersonId;
    this.__backAction = backAction;
  }
}
