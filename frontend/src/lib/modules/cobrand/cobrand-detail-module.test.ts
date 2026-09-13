import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { ICobrandService } from '$lib/services/cobrand/cobrand-service';
import { CobrandDetailModule } from './cobrand-detail-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

const mockCobrand = { id: 'c1', name: 'Test Firm' };
const mockAdmins = [{ id: 'a1', user: 'u1', cobrand: 'c1' }];
const mockManagers = [{ id: 'm1', cobrand: 'c1', property: 'p1' }];

const makeCobrandService = (overrides: Partial<ICobrandService> = {}) =>
  ({
    getCobrandById: vi.fn().mockResolvedValue(mockCobrand),
    getAdminsForCobrand: vi.fn().mockResolvedValue(mockAdmins),
    getPropertyManagersForCobrand: vi.fn().mockResolvedValue(mockManagers),
    updateCobrand: vi.fn().mockResolvedValue(undefined),
    deleteCobrand: vi.fn().mockResolvedValue(undefined),
    addAdmin: vi.fn().mockResolvedValue({ id: 'a2', user: 'u2', cobrand: 'c1' }),
    removeAdmin: vi.fn().mockResolvedValue(undefined),
    addPropertyManager: vi.fn().mockResolvedValue({ id: 'm2', cobrand: 'c1', property: 'p2' }),
    removePropertyManager: vi.fn().mockResolvedValue(undefined),
    ...overrides,
  }) as unknown as ICobrandService;

describe('CobrandDetailModule', () => {
  it('isAdd returns false', () => {
    const module = new CobrandDetailModule(makeCobrandService(), 'c1', vi.fn());
    expect(module.isAdd).toBe(false);
  });

  it('callBackAction calls the provided callback', () => {
    const backAction = vi.fn();
    const module = new CobrandDetailModule(makeCobrandService(), 'c1', backAction);
    module.callBackAction();
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('resolves cobrand from service', async () => {
    const module = new CobrandDetailModule(makeCobrandService(), 'c1', vi.fn());
    await expect(module.cobrandAsync).resolves.toEqual(mockCobrand);
  });

  it('resolves admins from service', async () => {
    const module = new CobrandDetailModule(makeCobrandService(), 'c1', vi.fn());
    await expect(module.adminsAsync).resolves.toEqual(mockAdmins);
  });

  it('resolves propertyManagers from service', async () => {
    const module = new CobrandDetailModule(makeCobrandService(), 'c1', vi.fn());
    await expect(module.propertyManagersAsync).resolves.toEqual(mockManagers);
  });

  it('getSaveCobrandAction calls service and navigates back', async () => {
    const backAction = vi.fn();
    const cobrandService = makeCobrandService();
    const module = new CobrandDetailModule(cobrandService, 'c1', backAction);
    await module.getSaveCobrandAction()(mockCobrand);
    expect(cobrandService.updateCobrand).toHaveBeenCalledWith('c1', 'Test Firm');
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('getDeleteCobrandAction calls service and navigates back', async () => {
    const backAction = vi.fn();
    const cobrandService = makeCobrandService();
    const module = new CobrandDetailModule(cobrandService, 'c1', backAction);
    await module.getDeleteCobrandAction()!('c1');
    expect(cobrandService.deleteCobrand).toHaveBeenCalledWith('c1');
    expect(backAction).toHaveBeenCalledOnce();
  });

  it('getAddAdminAction calls service and refreshes admins', async () => {
    const cobrandService = makeCobrandService();
    const module = new CobrandDetailModule(cobrandService, 'c1', vi.fn());
    await module.getAddAdminAction()('c1', 'u2');
    expect(cobrandService.addAdmin).toHaveBeenCalledWith('c1', 'u2');
    expect(cobrandService.getAdminsForCobrand).toHaveBeenCalledTimes(2);
  });

  it('getRemoveAdminAction calls service and refreshes admins', async () => {
    const cobrandService = makeCobrandService();
    const module = new CobrandDetailModule(cobrandService, 'c1', vi.fn());
    await module.getRemoveAdminAction()('a1');
    expect(cobrandService.removeAdmin).toHaveBeenCalledWith('a1');
    expect(cobrandService.getAdminsForCobrand).toHaveBeenCalledTimes(2);
  });

  it('getAddPropertyManagerAction calls service and refreshes managers', async () => {
    const cobrandService = makeCobrandService();
    const module = new CobrandDetailModule(cobrandService, 'c1', vi.fn());
    await module.getAddPropertyManagerAction()('c1', 'p2');
    expect(cobrandService.addPropertyManager).toHaveBeenCalledWith('c1', 'p2');
    expect(cobrandService.getPropertyManagersForCobrand).toHaveBeenCalledTimes(2);
  });

  it('getRemovePropertyManagerAction calls service and refreshes managers', async () => {
    const cobrandService = makeCobrandService();
    const module = new CobrandDetailModule(cobrandService, 'c1', vi.fn());
    await module.getRemovePropertyManagerAction()('m1');
    expect(cobrandService.removePropertyManager).toHaveBeenCalledWith('m1');
    expect(cobrandService.getPropertyManagersForCobrand).toHaveBeenCalledTimes(2);
  });
});
