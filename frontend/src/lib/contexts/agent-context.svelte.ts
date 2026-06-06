import { getContext, setContext } from 'svelte';

export class AgentContext {
  public selectedAgentId = $state<string>('');
}

export const setAgentContext = () => {
  const agentContext = new AgentContext();
  setContext<AgentContext>('agentContext', agentContext);
};

export const getAgentContext = () => getContext<AgentContext>('agentContext');
