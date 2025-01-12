import { getContext, setContext } from "svelte";

export class PersonContext {
  public selectedPersonId = $state<string>("");
}

export const setPersonContext = () => {
  const personContext = new PersonContext();
  setContext<PersonContext>("personContext", personContext)
};

export const getPersonContext = () => getContext<PersonContext>("personContext");