import type { IPropertyService } from '$lib/services/property/property-service';
import type { IPropertyDetailModule } from '$lib/modules/interfaces';
import type {
  IPropertyModel,
  IPersonPropertyOwnerModel,
  ITenantModel,
  IHouseholdModel,
} from '$lib/models/property-models';
import type { IPropertyAgentModel } from '$lib/models/agent-models';

export class PropertyDetailModule implements IPropertyDetailModule {
  private readonly __propertyService: IPropertyService;
  private readonly __propertyId: string;
  private readonly __backAction: () => void;

  public propertyAsync: Promise<IPropertyModel | null>;
  public personOwnersAsync: Promise<IPersonPropertyOwnerModel[]>;
  public tenantsAsync = $state<Promise<ITenantModel[]>>(Promise.resolve([]));
  public householdMembersAsync = $state<Promise<IHouseholdModel[]>>(Promise.resolve([]));
  public propertyAgentsAsync: Promise<IPropertyAgentModel[]>;

  get isAdd() {
    return false;
  }

  public callBackAction = () => this.__backAction();

  private refreshTenants() {
    this.tenantsAsync = this.__propertyService
      .getTenantsForProperty(this.__propertyId)
      .catch((ex) => {
        alert(ex);
        return [];
      });
  }

  private refreshHouseholdMembers() {
    this.householdMembersAsync = this.__propertyService
      .getHouseholdMembersForProperty(this.__propertyId)
      .catch((ex) => {
        alert(ex);
        return [];
      });
  }

  public getSavePropertyAction = () => async (property: IPropertyModel) => {
    try {
      await this.__propertyService.updateProperty(property.id, property.address);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getDeletePropertyAction = () => async (id: string) => {
    try {
      await this.__propertyService.deleteProperty(id);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getAddTenantAction = () => async (personId: string, propertyId: string) => {
    try {
      await this.__propertyService.addTenant(personId, propertyId);
      this.refreshTenants();
    } catch (ex) {
      alert(ex);
    }
  };

  public getRemoveTenantAction = () => async (id: string) => {
    try {
      await this.__propertyService.removeTenant(id);
      this.refreshTenants();
    } catch (ex) {
      alert(ex);
    }
  };

  public getAddHouseholdMemberAction = () => async (personId: string, propertyId: string) => {
    try {
      await this.__propertyService.addHouseholdMember(personId, propertyId);
      this.refreshHouseholdMembers();
    } catch (ex) {
      alert(ex);
    }
  };

  public getRemoveHouseholdMemberAction = () => async (id: string) => {
    try {
      await this.__propertyService.removeHouseholdMember(id);
      this.refreshHouseholdMembers();
    } catch (ex) {
      alert(ex);
    }
  };

  constructor(propertyService: IPropertyService, propertyId: string, backAction: () => void) {
    this.__propertyService = propertyService;
    this.__propertyId = propertyId;
    this.__backAction = backAction;

    this.propertyAsync = this.__propertyService.getPropertyById(propertyId).catch((ex) => {
      alert(ex);
      return null;
    });

    this.personOwnersAsync = this.__propertyService
      .getPersonOwnersForProperty(propertyId)
      .catch((ex) => {
        alert(ex);
        return [];
      });

    this.propertyAgentsAsync = this.__propertyService
      .getPropertyAgentsForProperty(propertyId)
      .catch((ex) => {
        alert(ex);
        return [];
      });

    this.refreshTenants();
    this.refreshHouseholdMembers();
  }
}
