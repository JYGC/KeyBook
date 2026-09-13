import type { IAgentService } from '$lib/services/agent/agent-service';
import type { IAgentDetailModule } from '$lib/modules/interfaces';
import type { IAgentModel, IPropertyAgentModel } from '$lib/models/agent-models';

export class AgentDetailModule implements IAgentDetailModule {
  private readonly __agentService: IAgentService;
  private readonly __agentId: string;
  private readonly __backAction: () => void;

  public agentAsync: Promise<IAgentModel | null>;
  public propertyAgentsAsync = $state<Promise<IPropertyAgentModel[]>>(Promise.resolve([]));

  get isAdd() {
    return false;
  }

  public callBackAction = () => this.__backAction();

  private refreshPropertyAgents() {
    this.propertyAgentsAsync = this.__agentService
      .getPropertyAgentsForAgent(this.__agentId)
      .catch((error) => {
        alert(error);
        return [];
      });
  }

  public getSaveAgentAction = () => null;

  public getDeleteAgentAction = () => async (id: string) => {
    try {
      await this.__agentService.deleteAgent(id);
      this.__backAction();
    } catch (error) {
      alert(error);
    }
  };

  public getAddPropertyAgentAction = () => async (agentId: string, propertyId: string) => {
    try {
      await this.__agentService.addPropertyAgent(agentId, propertyId);
      this.refreshPropertyAgents();
    } catch (error) {
      alert(error);
    }
  };

  public getRemovePropertyAgentAction = () => async (id: string) => {
    try {
      await this.__agentService.removePropertyAgent(id);
      this.refreshPropertyAgents();
    } catch (error) {
      alert(error);
    }
  };

  constructor(agentService: IAgentService, agentId: string, backAction: () => void) {
    this.__agentService = agentService;
    this.__agentId = agentId;
    this.__backAction = backAction;

    this.agentAsync = this.__agentService
      .getAgentById(agentId)
      .catch((error) => {
        alert(error);
        return null;
      });

    this.refreshPropertyAgents();
  }
}
