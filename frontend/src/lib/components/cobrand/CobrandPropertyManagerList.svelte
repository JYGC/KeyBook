<script lang="ts">
  import type { ICobrandDetailModule } from '$lib/modules/interfaces';
  import { Button, ButtonSet, DataTable, TextInput, Tile } from 'carbon-components-svelte';

  let {
    cobrandId,
    cobrandDetailModule,
  } = $props<{
    cobrandId: string;
    cobrandDetailModule: ICobrandDetailModule;
  }>();

  const addManagerAction = cobrandDetailModule.getAddPropertyManagerAction();
  const removeManagerAction = cobrandDetailModule.getRemovePropertyManagerAction();

  let newPropertyId = $state('');
</script>

<h3>Property Managers</h3>

{#await cobrandDetailModule.propertyManagersAsync}
  <Tile>...getting property managers</Tile>
{:then managers}
  <DataTable
    headers={[
      { key: 'property', value: 'Property ID' },
      { key: 'id', empty: true },
    ]}
    rows={managers}
  >
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === 'id'}
        <ButtonSet>
          <Button kind="danger" onclick={() => removeManagerAction(cell.value)}>Remove</Button>
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
    addManagerAction(cobrandId, newPropertyId);
    newPropertyId = '';
  }}
>
  Add Property Manager
</Button>
