<script lang="ts">
  import { goto } from '$app/navigation';
  import type { INewPersonListModule } from '$lib/modules/interfaces';
  import { Button, ButtonSet, DataTable, Tile } from 'carbon-components-svelte';

  let {
    personListModule,
  } = $props<{
    personListModule: INewPersonListModule;
  }>();

  const goToDetail = (personId: string) => {
    goto(`/user/persons/detail?id=${personId}`);
  };
</script>

{#await personListModule.personListAsync}
  <Tile>...getting persons</Tile>
{:then personList}
  <DataTable
    headers={[
      { key: 'name', value: 'Name' },
      { key: 'DOB', value: 'Date of Birth' },
      { key: 'id', empty: true },
    ]}
    rows={personList}
  >
    <strong slot="title">Persons</strong>
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
