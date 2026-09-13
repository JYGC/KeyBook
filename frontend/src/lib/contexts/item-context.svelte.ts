import { getContext, setContext } from 'svelte';

export class ItemContext {
  public selectedItemId = $state<string>('');
}

export const setItemContext = () => {
  const itemContext = new ItemContext();
  setContext<ItemContext>('itemContext', itemContext);
};

export const getItemContext = () => getContext<ItemContext>('itemContext');
