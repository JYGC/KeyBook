import type { IPersonModel } from "$lib/models/person-models";
import type { IPropertyModel, IPersonPropertyOwnerModel, ITenantModel, IHouseholdModel } from "$lib/models/property-models";
import type { IItemModel, IEntryDeviceModel } from "$lib/models/item-models";
import type { ICobrandModel, ICobrandAdminModel, ICobrandPropertyManagerModel } from "$lib/models/cobrand-models";
import type { IAgentModel, IPropertyAgentModel } from "$lib/models/agent-models";

export interface INewPersonListModule {
  personListAsync: Promise<IPersonModel[]>;
}

export interface INewPersonEditorModule {
  personAsync: Promise<IPersonModel | null>;
  get isAdd(): boolean;
  getSavePersonAction: () => (person: IPersonModel) => Promise<void>;
  getDeletePersonAction: () => ((id: string) => Promise<void>) | null;
  callBackAction: () => void;
}

export interface IPersonDetailModule extends INewPersonEditorModule {
  rolesAsync: Promise<string[]>;
}

export interface INewPropertyListModule {
  propertyListAsync: Promise<IPropertyModel[]>;
}

export interface INewPropertyEditorModule {
  propertyAsync: Promise<IPropertyModel | null>;
  get isAdd(): boolean;
  getSavePropertyAction: () => (property: IPropertyModel) => Promise<void>;
  getDeletePropertyAction: () => ((id: string) => Promise<void>) | null;
  callBackAction: () => void;
}

export interface IPropertyDetailModule extends INewPropertyEditorModule {
  personOwnersAsync: Promise<IPersonPropertyOwnerModel[]>;
  tenantsAsync: Promise<ITenantModel[]>;
  householdMembersAsync: Promise<IHouseholdModel[]>;
  propertyAgentsAsync: Promise<IPropertyAgentModel[]>;
  getAddTenantAction: () => (personId: string, propertyId: string) => Promise<void>;
  getRemoveTenantAction: () => (id: string) => Promise<void>;
  getAddHouseholdMemberAction: () => (personId: string, propertyId: string) => Promise<void>;
  getRemoveHouseholdMemberAction: () => (id: string) => Promise<void>;
}

export interface IItemListModule {
  itemListAsync: Promise<IItemModel[]>;
}

export interface IItemEditorModule {
  itemAsync: Promise<IItemModel | null>;
  entryDeviceAsync: Promise<IEntryDeviceModel | null>;
  get isAdd(): boolean;
  getSaveItemAction: () => (item: IItemModel) => Promise<void>;
  getDeleteItemAction: () => ((id: string) => Promise<void>) | null;
  callBackAction: () => void;
}

export interface ICobrandListModule {
  cobrandListAsync: Promise<ICobrandModel[]>;
}

export interface ICobrandEditorModule {
  cobrandAsync: Promise<ICobrandModel | null>;
  get isAdd(): boolean;
  getSaveCobrandAction: () => (cobrand: ICobrandModel) => Promise<void>;
  getDeleteCobrandAction: () => ((id: string) => Promise<void>) | null;
  callBackAction: () => void;
}

export interface ICobrandDetailModule extends ICobrandEditorModule {
  adminsAsync: Promise<ICobrandAdminModel[]>;
  propertyManagersAsync: Promise<ICobrandPropertyManagerModel[]>;
  getAddAdminAction: () => (cobrandId: string, userId: string) => Promise<void>;
  getRemoveAdminAction: () => (id: string) => Promise<void>;
  getAddPropertyManagerAction: () => (cobrandId: string, propertyId: string) => Promise<void>;
  getRemovePropertyManagerAction: () => (id: string) => Promise<void>;
}

export interface IAgentListModule {
  agentListAsync: Promise<IAgentModel[]>;
}

export interface IAgentEditorModule {
  agentAsync: Promise<IAgentModel | null>;
  get isAdd(): boolean;
  getSaveAgentAction: () => ((agent: IAgentModel) => Promise<void>) | null;
  getDeleteAgentAction: () => ((id: string) => Promise<void>) | null;
  callBackAction: () => void;
}

export interface IAgentDetailModule extends IAgentEditorModule {
  propertyAgentsAsync: Promise<IPropertyAgentModel[]>;
  getAddPropertyAgentAction: () => (agentId: string, propertyId: string) => Promise<void>;
  getRemovePropertyAgentAction: () => (id: string) => Promise<void>;
}

export interface ILoginModule {
  email: string;
  password: string;
  callApi: () => Promise<string>;
}

export interface IRegisterModule {
  name: string;
  email: string;
  password: string;
  passwordConfirm: string;
  callApi: () => Promise<boolean>;
  get error(): string;
}
