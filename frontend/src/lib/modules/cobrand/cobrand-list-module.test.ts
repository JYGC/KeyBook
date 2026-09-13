import { describe, it, expect, vi, beforeAll } from 'vitest';
import type { ICobrandService } from '$lib/services/cobrand/cobrand-service';
import { CobrandListModule } from './cobrand-list-module.svelte';

beforeAll(() => {
  vi.stubGlobal('alert', vi.fn());
});

describe('CobrandListModule', () => {
  it('returns cobrands from service', async () => {
    const cobrands = [
      { id: '1', name: 'Acme Real Estate' },
      { id: '2', name: 'Global Properties' },
    ];
    const cobrandService = { getAllCobrands: vi.fn().mockResolvedValue(cobrands) } as unknown as ICobrandService;

    const module = new CobrandListModule(cobrandService);
    await expect(module.cobrandListAsync).resolves.toEqual(cobrands);
    expect(cobrandService.getAllCobrands).toHaveBeenCalledOnce();
  });

  it('returns empty array when service throws', async () => {
    const cobrandService = {
      getAllCobrands: vi.fn().mockRejectedValue(new Error('network error')),
    } as unknown as ICobrandService;

    const module = new CobrandListModule(cobrandService);
    await expect(module.cobrandListAsync).resolves.toEqual([]);
  });
});
