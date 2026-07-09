import type { IAgentService } from '$lib/services/agent/agent-service';
import type { IAgentEditorModule } from '$lib/modules/interfaces';
import type { IAgentModel } from '$lib/models/agent-models';

export class AgentAddModule implements IAgentEditorModule {
  private readonly __agentService: IAgentService;
  private readonly __backAction: () => void;

  public agentAsync = $state<Promise<IAgentModel | null>>(
    Promise.resolve({ id: '', person: '', cobrand: '' }),
  );

  get isAdd() {
    return true;
  }

  public callBackAction = () => this.__backAction();

  public getSaveAgentAction = () => async (agent: IAgentModel) => {
    try {
      await this.__agentService.createAgent(agent.person, agent.cobrand);
      this.__backAction();
    } catch (ex) {
      alert(ex);
    }
  };

  public getDeleteAgentAction = () => null;

  constructor(agentService: IAgentService, backAction: () => void) {
    this.__agentService = agentService;
    this.__backAction = backAction;
  }
}
