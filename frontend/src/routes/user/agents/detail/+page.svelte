<script lang="ts">
  import { getBackendClient } from '$lib/api/backend-client';
  import AgentEditor from '$lib/components/agent/AgentEditor.svelte';
  import PropertyAgentList from '$lib/components/agent/PropertyAgentList.svelte';
  import ConfirmButtonAndDialog from '$lib/components/shared/ConfirmButtonAndDialog.svelte';
  import { AgentRepository } from '$lib/repositories/agent/agent-repository';
  import { PropertyAgentRepository } from '$lib/repositories/agent/property-agent-repository';
  import { AgentService } from '$lib/services/agent/agent-service';
  import { AgentDetailModule } from '$lib/modules/agent/agent-detail-module.svelte';
  import { Button } from 'carbon-components-svelte';

  let { data } = $props();
  const agentId: string = data.agentId;

  const goBack = () => {
    history.back();
  };

  const pb = getBackendClient();
  const agentRepo = new AgentRepository(pb);
  const propertyAgentRepo = new PropertyAgentRepository(pb);
  const agentService = new AgentService(agentRepo, propertyAgentRepo);
  const agentDetailModule = new AgentDetailModule(agentService, agentId, goBack);
</script>

<Button onclick={goBack}>Back</Button>

<AgentEditor agentEditorModule={agentDetailModule}>
  {#snippet deleteButton(deleteActionButtonClick: () => void)}
    <ConfirmButtonAndDialog
      submitAction={deleteActionButtonClick}
      buttonText="Delete Agent"
      bodyMessage="Are you sure you want to delete this agent?"
    />
  {/snippet}
</AgentEditor>

<br />

<PropertyAgentList {agentId} {agentDetailModule} />
