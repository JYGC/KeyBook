<script lang="ts">
  import { goto } from '$app/navigation';
  import type { IItemListModule } from '$lib/modules/interfaces';
  import { Button, ButtonSet, DataTable, Tile } from 'carbon-components-svelte';

  let {
    itemListModule,
  } = $props<{
    itemListModule: IItemListModule;
  }>();

  const goToEditItem = (itemId: string) => {
    goto(`/user/items/edit?id=${itemId}`);
  };
</script>

{#await itemListModule.itemListAsync}
  <Tile>...getting items</Tile>
{:then itemList}
  <DataTable
    headers={[
      { key: 'name', value: 'Name' },
      { key: 'description', value: 'Description' },
      { key: 'id', empty: true },
    ]}
    rows={itemList}
  >
    <strong slot="title">Items</strong>
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === 'id'}
        <ButtonSet>
          <Button onclick={() => goToEditItem(cell.value)}>Edit</Button>
        </ButtonSet>
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}
