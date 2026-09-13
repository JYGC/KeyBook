<script lang="ts">
  import type { IAgentDetailModule } from '$lib/modules/interfaces';
  import { Button, ButtonSet, DataTable, TextInput, Tile } from 'carbon-components-svelte';

  let {
    agentId,
    agentDetailModule,
  } = $props<{
    agentId: string;
    agentDetailModule: IAgentDetailModule;
  }>();

  const addPropertyAgentAction = agentDetailModule.getAddPropertyAgentAction();
  const removePropertyAgentAction = agentDetailModule.getRemovePropertyAgentAction();

  let newPropertyId = $state('');
</script>

<h3>Property Assignments</h3>

{#await agentDetailModule.propertyAgentsAsync}
  <Tile>...getting property assignments</Tile>
{:then propertyAgents}
  <DataTable
    headers={[
      { key: 'property', value: 'Property ID' },
      { key: 'id', empty: true },
    ]}
    rows={propertyAgents}
  >
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === 'id'}
        <ButtonSet>
          <Button kind="danger" onclick={() => removePropertyAgentAction(cell.value)}>Remove</Button>
        </ButtonSet>
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}

<br />
<TextInput labelText="Property ID" bind:value={newPropertyId} />
<br />
<Button
  onclick={() => {
    addPropertyAgentAction(agentId, newPropertyId);
    newPropertyId = '';
  }}
>
  Add Property Assignment
</Button>
