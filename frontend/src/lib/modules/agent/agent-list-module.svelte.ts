import type { IAgentService } from '$lib/services/agent/agent-service';
import type { IAgentListModule } from '$lib/modules/interfaces';
import type { IAgentModel } from '$lib/models/agent-models';

export class AgentListModule implements IAgentListModule {
  private readonly __agentService: IAgentService;

  public agentListAsync = $derived.by<Promise<IAgentModel[]>>(async () => {
    try {
      return await this.__agentService.getAllAgents();
    } catch (ex) {
      alert(ex);
      return [];
    }
  });

  constructor(agentService: IAgentService) {
    this.__agentService = agentService;
  }
}
