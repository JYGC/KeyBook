import { getContext, setContext } from 'svelte';

export class CobrandContext {
  public selectedCobrandId = $state<string>('');
}

export const setCobrandContext = () => {
  const cobrandContext = new CobrandContext();
  setContext<CobrandContext>('cobrandContext', cobrandContext);
};

export const getCobrandContext = () => getContext<CobrandContext>('cobrandContext');
