<script lang="ts">
  import type { IPropertyDetailModule } from '$lib/modules/interfaces';
  import { Button, ButtonSet, DataTable, TextInput, Tile } from 'carbon-components-svelte';

  let {
    propertyId,
    propertyDetailModule,
  } = $props<{
    propertyId: string;
    propertyDetailModule: IPropertyDetailModule;
  }>();

  const addTenantAction = propertyDetailModule.getAddTenantAction();
  const removeTenantAction = propertyDetailModule.getRemoveTenantAction();

  let newPersonId = $state('');
</script>

<h3>Tenants</h3>

{#await propertyDetailModule.tenantsAsync}
  <Tile>...getting tenants</Tile>
{:then tenants}
  <DataTable
    headers={[
      { key: 'person', value: 'Person ID' },
      { key: 'id', empty: true },
    ]}
    rows={tenants}
  >
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === 'id'}
        <ButtonSet>
          <Button kind="danger" onclick={() => removeTenantAction(cell.value)}>Remove</Button>
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
<TextInput labelText="Person ID" bind:value={newPersonId} />
<br />
<Button
  onclick={() => {
    addTenantAction(newPersonId, propertyId);
    newPersonId = '';
  }}
>
  Add Tenant
</Button>
