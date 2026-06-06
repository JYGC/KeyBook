<script lang="ts">
  import { goto } from '$app/navigation';
  import { getBackendClient } from '$lib/api/backend-client';
  import AgentList from '$lib/components/agent/AgentList.svelte';
  import { AgentRepository } from '$lib/repositories/agent/agent-repository';
  import { PropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';
  import { AgentService } from '$lib/services/agent/agent-service';
  import { AgentListModule } from '$lib/modules/agent/agent-list-module.svelte';
  import { Button } from 'carbon-components-svelte';

  const pb = getBackendClient();
  const agentRepo = new AgentRepository(pb);
  const propertyAgentRepo = new PropertyAgentRepository(pb);
  const agentService = new AgentService(agentRepo, propertyAgentRepo);
  const agentListModule = new AgentListModule(agentService);

  const gotoAddAgent = () => {
    goto('/user/agents/add');
  };
</script>

<Button onclick={gotoAddAgent}>Add Agent</Button>
<AgentList {agentListModule} />
