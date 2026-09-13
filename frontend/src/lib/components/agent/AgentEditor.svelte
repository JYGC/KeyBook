<script lang="ts">
  import type { IAgentEditorModule } from '$lib/modules/interfaces';
  import { Button, TextInput, Tile } from 'carbon-components-svelte';
  import type { Snippet } from 'svelte';

  let {
    deleteButton,
    agentEditorModule,
  } = $props<{
    deleteButton?: Snippet<[() => void]>;
    agentEditorModule: IAgentEditorModule;
  }>();

  const saveAction = agentEditorModule.getSaveAgentAction();
  const deleteActionButtonClick = agentEditorModule.getDeleteAgentAction();
</script>

{#await agentEditorModule.agentAsync}
  <Tile>...getting agent details</Tile>
{:then agent}
  {#if agent === null}
    <p>Agent not found.</p>
  {:else}
    <TextInput labelText="Person ID" bind:value={agent.person} disabled={!agentEditorModule.isAdd} />
    <br />
    <TextInput labelText="Cobrand ID" bind:value={agent.cobrand} disabled={!agentEditorModule.isAdd} />
    <br />
    <br />
    {#if saveAction !== null}
      <Button onclick={() => saveAction(agent)}>Save</Button>
    {/if}
    {#if deleteButton !== undefined && deleteActionButtonClick !== undefined}
      {@render deleteButton(() => deleteActionButtonClick(agent.id))}
    {/if}
  {/if}
{:catch error}
  {error}
{/await}
