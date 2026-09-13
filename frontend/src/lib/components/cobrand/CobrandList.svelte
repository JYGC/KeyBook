<script lang="ts">
  import { goto } from '$app/navigation';
  import type { ICobrandListModule } from '$lib/modules/interfaces';
  import { Button, ButtonSet, DataTable, Tile } from 'carbon-components-svelte';

  let {
    cobrandListModule,
  } = $props<{
    cobrandListModule: ICobrandListModule;
  }>();

  const goToDetail = (cobrandId: string) => {
    goto(`/user/cobrands/detail?id=${cobrandId}`);
  };
</script>

{#await cobrandListModule.cobrandListAsync}
  <Tile>...getting cobrands</Tile>
{:then cobrandList}
  <DataTable
    headers={[
      { key: 'name', value: 'Name' },
      { key: 'id', empty: true },
    ]}
    rows={cobrandList}
  >
    <strong slot="title">Cobrands</strong>
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === 'id'}
        <ButtonSet>
          <Button onclick={() => goToDetail(cell.value)}>Edit</Button>
        </ButtonSet>
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}
