<script lang="ts">
  import type { IPropertyDetailModule } from '$lib/modules/interfaces';
  import { DataTable, Tile } from 'carbon-components-svelte';

  let {
    propertyDetailModule,
  } = $props<{
    propertyDetailModule: IPropertyDetailModule;
  }>();
</script>

<h3>Person Owners</h3>

{#await propertyDetailModule.personOwnersAsync}
  <Tile>...getting person owners</Tile>
{:then personOwners}
  <DataTable
    headers={[
      { key: 'person', value: 'Person ID' },
      { key: 'propertyOwner', value: 'Property Owner ID' },
    ]}
    rows={personOwners}
  >
    <svelte:fragment slot="cell" let:cell>
      {cell.value}
    </svelte:fragment>
  </DataTable>
{:catch error}
  {error}
{/await}
