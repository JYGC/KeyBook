<script lang="ts">
  import { goto } from '$app/navigation';
  import type { IAgentListModule } from '$lib/modules/interfaces';
  import { Button, ButtonSet, DataTable, Tile } from 'carbon-components-svelte';

  let {
    agentListModule,
  } = $props<{
    agentListModule: IAgentListModule;
  }>();

  const goToDetail = (agentId: string) => {
    goto(`/user/agents/detail?id=${agentId}`);
  };
</script>

{#await agentListModule.agentListAsync}
  <Tile>...getting agents</Tile>
{:then agentList}
  <DataTable
    headers={[
      { key: 'person', value: 'Person ID' },
      { key: 'cobrand', value: 'Cobrand ID' },
      { key: 'id', empty: true },
    ]}
    rows={agentList}
  >
    <strong slot="title">Agents</strong>
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === 'id'}
        <ButtonSet>
          <Button onclick={() => goToDetail(cell.value)}>Detail</Button>
        </ButtonSet>
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}
