import { getContext, setContext } from "svelte";

export class PropertyContext {
  public selectedPropertyId = $state<string>("");
}

export const setPropertyContext = () => {
  const propertyContext = new PropertyContext();
  setContext<PropertyContext>("propertyContext", propertyContext)
};

export const getPropertyContext = () => getContext<PropertyContext>("propertyContext");