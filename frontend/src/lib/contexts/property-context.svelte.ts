import type { IPropertyContext } from "$lib/interfaces";
import { getContext, setContext } from "svelte";

export class PropertyContext implements IPropertyContext {
  public selectedPropertyId = $state<string>("");
}

export const setPropertyContext = () => {
  const propertyContext = new PropertyContext();
  setContext<IPropertyContext>("propertyContext", propertyContext)
};

export const getPropertyContext = () => getContext<IPropertyContext>("propertyContext");