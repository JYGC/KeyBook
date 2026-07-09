import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { IPropertyService } from '$lib/services/property/property-service';
import { PropertyDetailModule } from './property-detail-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

const mockProperty = { id: 'pr1', address: '1 Main St' };
const mockPersonOwners = [{ id: 'ppo1', person: 'p1', propertyOwner: 'po1' }];
const mockTenants = [{ id: 't1', person: 'p2', property: 'pr1' }];
const mockHousehold = [{ id: 'h1', person: 'p3', property: 'pr1' }];
const mockPropertyAgents = [{ id: 'pa1', agent: 'ag1', property: 'pr1' }];

const makeSvc = (overrides: Partial<IPropertyService> = {}) =>
  ({
    getPropertyById: vi.fn().mockResolvedValue(mockProperty),
    getPersonOwnersForProperty: vi.fn().mockResolvedValue(mockPersonOwners),
    getTenantsForProperty: vi.fn().mockResolvedValue(mockTenants),
    getHouseholdMembersForProperty: vi.fn().mockResolvedValue(mockHousehold),
    getPropertyAgentsForProperty: vi.fn().mockResolvedValue(mockPropertyAgents),
    updateProperty: vi.fn().mockResolvedValue(mockProperty),
    deleteProperty: vi.fn().mockResolvedValue(undefined),
    addTenant: vi.fn().mockResolvedValue({ id: 't2', person: 'p4', property: 'pr1' }),
    removeTenant: vi.fn().mockResolvedValue(undefined),
    addHouseholdMember: vi.fn().mockResolvedValue({ id: 'h2', person: 'p5', property: 'pr1' }),
    removeHouseholdMember: vi.fn().mockResolvedValue(undefined),
    ...overrides,
  }) as unknown as IPropertyService;

describe('PropertyDetailModule', () => {
  it('isAdd returns false', () => {
    const module = new PropertyDetailModule(makeSvc(), 'pr1', vi.fn());
    expect(module.isAdd).toBe(false);
  });

  it('callBackAction calls the provided callback', () => {
    const back = vi.fn();
    const module = new PropertyDetailModule(makeSvc(), 'pr1', back);
    module.callBackAction();
    expect(back).toHaveBeenCalledOnce();
  });

  it('resolves property from service', async () => {
    const module = new PropertyDetailModule(makeSvc(), 'pr1', vi.fn());
    await expect(module.propertyAsync).resolves.toEqual(mockProperty);
  });

  it('resolves personOwners from service', async () => {
    const module = new PropertyDetailModule(makeSvc(), 'pr1', vi.fn());
    await expect(module.personOwnersAsync).resolves.toEqual(mockPersonOwners);
  });

  it('resolves tenants from service', async () => {
    const module = new PropertyDetailModule(makeSvc(), 'pr1', vi.fn());
    await expect(module.tenantsAsync).resolves.toEqual(mockTenants);
  });

  it('resolves householdMembers from service', async () => {
    const module = new PropertyDetailModule(makeSvc(), 'pr1', vi.fn());
    await expect(module.householdMembersAsync).resolves.toEqual(mockHousehold);
  });

  it('resolves propertyAgents from service', async () => {
    const module = new PropertyDetailModule(makeSvc(), 'pr1', vi.fn());
    await expect(module.propertyAgentsAsync).resolves.toEqual(mockPropertyAgents);
  });

  it('getSavePropertyAction calls updateProperty and navigates back', async () => {
    const back = vi.fn();
    const svc = makeSvc();
    const module = new PropertyDetailModule(svc, 'pr1', back);
    await module.getSavePropertyAction()(mockProperty);
    expect(svc.updateProperty).toHaveBeenCalledWith('pr1', '1 Main St');
    expect(back).toHaveBeenCalledOnce();
  });

  it('getDeletePropertyAction calls deleteProperty and navigates back', async () => {
    const back = vi.fn();
    const svc = makeSvc();
    const module = new PropertyDetailModule(svc, 'pr1', back);
    await module.getDeletePropertyAction()!('pr1');
    expect(svc.deleteProperty).toHaveBeenCalledWith('pr1');
    expect(back).toHaveBeenCalledOnce();
  });

  it('getAddTenantAction calls service and refreshes tenants', async () => {
    const svc = makeSvc();
    const module = new PropertyDetailModule(svc, 'pr1', vi.fn());
    await module.getAddTenantAction()('p4', 'pr1');
    expect(svc.addTenant).toHaveBeenCalledWith('p4', 'pr1');
    expect(svc.getTenantsForProperty).toHaveBeenCalledTimes(2);
  });

  it('getRemoveTenantAction calls service and refreshes tenants', async () => {
    const svc = makeSvc();
    const module = new PropertyDetailModule(svc, 'pr1', vi.fn());
    await module.getRemoveTenantAction()('t1');
    expect(svc.removeTenant).toHaveBeenCalledWith('t1');
    expect(svc.getTenantsForProperty).toHaveBeenCalledTimes(2);
  });

  it('getAddHouseholdMemberAction calls service and refreshes householdMembers', async () => {
    const svc = makeSvc();
    const module = new PropertyDetailModule(svc, 'pr1', vi.fn());
    await module.getAddHouseholdMemberAction()('p5', 'pr1');
    expect(svc.addHouseholdMember).toHaveBeenCalledWith('p5', 'pr1');
    expect(svc.getHouseholdMembersForProperty).toHaveBeenCalledTimes(2);
  });

  it('getRemoveHouseholdMemberAction calls service and refreshes householdMembers', async () => {
    const svc = makeSvc();
    const module = new PropertyDetailModule(svc, 'pr1', vi.fn());
    await module.getRemoveHouseholdMemberAction()('h1');
    expect(svc.removeHouseholdMember).toHaveBeenCalledWith('h1');
    expect(svc.getHouseholdMembersForProperty).toHaveBeenCalledTimes(2);
  });
});
