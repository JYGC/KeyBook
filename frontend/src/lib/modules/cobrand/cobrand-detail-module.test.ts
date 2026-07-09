import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { ICobrandService } from '$lib/services/cobrand/cobrand-service';
import { CobrandDetailModule } from './cobrand-detail-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

const mockCobrand = { id: 'c1', name: 'Test Firm' };
const mockAdmins = [{ id: 'a1', user: 'u1', cobrand: 'c1' }];
const mockManagers = [{ id: 'm1', cobrand: 'c1', property: 'p1' }];

const makeSvc = (overrides: Partial<ICobrandService> = {}) =>
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
    const module = new CobrandDetailModule(makeSvc(), 'c1', vi.fn());
    expect(module.isAdd).toBe(false);
  });

  it('callBackAction calls the provided callback', () => {
    const back = vi.fn();
    const module = new CobrandDetailModule(makeSvc(), 'c1', back);
    module.callBackAction();
    expect(back).toHaveBeenCalledOnce();
  });

  it('resolves cobrand from service', async () => {
    const module = new CobrandDetailModule(makeSvc(), 'c1', vi.fn());
    await expect(module.cobrandAsync).resolves.toEqual(mockCobrand);
  });

  it('resolves admins from service', async () => {
    const module = new CobrandDetailModule(makeSvc(), 'c1', vi.fn());
    await expect(module.adminsAsync).resolves.toEqual(mockAdmins);
  });

  it('resolves propertyManagers from service', async () => {
    const module = new CobrandDetailModule(makeSvc(), 'c1', vi.fn());
    await expect(module.propertyManagersAsync).resolves.toEqual(mockManagers);
  });

  it('getSaveCobrandAction calls service and navigates back', async () => {
    const back = vi.fn();
    const svc = makeSvc();
    const module = new CobrandDetailModule(svc, 'c1', back);
    await module.getSaveCobrandAction()(mockCobrand);
    expect(svc.updateCobrand).toHaveBeenCalledWith('c1', 'Test Firm');
    expect(back).toHaveBeenCalledOnce();
  });

  it('getDeleteCobrandAction calls service and navigates back', async () => {
    const back = vi.fn();
    const svc = makeSvc();
    const module = new CobrandDetailModule(svc, 'c1', back);
    await module.getDeleteCobrandAction()!('c1');
    expect(svc.deleteCobrand).toHaveBeenCalledWith('c1');
    expect(back).toHaveBeenCalledOnce();
  });

  it('getAddAdminAction calls service and refreshes admins', async () => {
    const svc = makeSvc();
    const module = new CobrandDetailModule(svc, 'c1', vi.fn());
    await module.getAddAdminAction()('c1', 'u2');
    expect(svc.addAdmin).toHaveBeenCalledWith('c1', 'u2');
    expect(svc.getAdminsForCobrand).toHaveBeenCalledTimes(2);
  });

  it('getRemoveAdminAction calls service and refreshes admins', async () => {
    const svc = makeSvc();
    const module = new CobrandDetailModule(svc, 'c1', vi.fn());
    await module.getRemoveAdminAction()('a1');
    expect(svc.removeAdmin).toHaveBeenCalledWith('a1');
    expect(svc.getAdminsForCobrand).toHaveBeenCalledTimes(2);
  });

  it('getAddPropertyManagerAction calls service and refreshes managers', async () => {
    const svc = makeSvc();
    const module = new CobrandDetailModule(svc, 'c1', vi.fn());
    await module.getAddPropertyManagerAction()('c1', 'p2');
    expect(svc.addPropertyManager).toHaveBeenCalledWith('c1', 'p2');
    expect(svc.getPropertyManagersForCobrand).toHaveBeenCalledTimes(2);
  });

  it('getRemovePropertyManagerAction calls service and refreshes managers', async () => {
    const svc = makeSvc();
    const module = new CobrandDetailModule(svc, 'c1', vi.fn());
    await module.getRemovePropertyManagerAction()('m1');
    expect(svc.removePropertyManager).toHaveBeenCalledWith('m1');
    expect(svc.getPropertyManagersForCobrand).toHaveBeenCalledTimes(2);
  });
});
