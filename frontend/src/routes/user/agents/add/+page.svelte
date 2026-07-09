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

  const pb = getBackendClient();
  const agentRepo = new AgentRepository(pb);
  const propertyAgentRepo = new PropertyAgentRepository(pb);
  const agentService = new AgentService(agentRepo, propertyAgentRepo);
  const agentAddModule = new AgentAddModule(agentService, goBack);
</script>

<Button onclick={goBack}>Back</Button>
<AgentEditor agentEditorModule={agentAddModule} />
