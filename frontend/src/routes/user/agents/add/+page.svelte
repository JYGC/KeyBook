<script lang="ts">
  import { getBackendClient } from '$lib/api/backend-client';
  import AgentEditor from '$lib/components/agent/AgentEditor.svelte';
  import { AgentRepository } from '$lib/repositories/agent/agent-repository';
  import { PropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';
  import { AgentService } from '$lib/services/agent/agent-service';
  import { AgentAddModule } from '$lib/modules/agent/agent-add-module.svelte';
  import { Button } from 'carbon-components-svelte';

  const goBack = () => {
    history.back();
  };

  const backendClient = getBackendClient();
  const agentRepository = new AgentRepository(backendClient);
  const propertyAgentRepository = new PropertyAgentRepository(backendClient);
  const agentService = new AgentService(agentRepository, propertyAgentRepository);
  const agentAddModule = new AgentAddModule(agentService, goBack);
</script>

<Button onclick={goBack}>Back</Button>
<AgentEditor agentEditorModule={agentAddModule} />
