<script lang="ts">
  import { goto } from '$app/navigation';
  import type { INewPropertyListModule } from '$lib/modules/interfaces';
  import { Button, ButtonSet, DataTable, Tile } from 'carbon-components-svelte';

  let {
    propertyListModule,
  } = $props<{
    propertyListModule: INewPropertyListModule;
  }>();

  const goToDetail = (propertyId: string) => {
    goto(`/user/properties/detail?id=${propertyId}`);
  };
</script>

{#await propertyListModule.propertyListAsync}
  <Tile>...getting properties</Tile>
{:then propertyList}
  <DataTable
    headers={[
      { key: 'address', value: 'Property Address' },
      { key: 'id', empty: true },
    ]}
    rows={propertyList}
  >
    <strong slot="title">Properties</strong>
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
