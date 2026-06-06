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

  const addAdminAction = cobrandDetailModule.getAddAdminAction();
  const removeAdminAction = cobrandDetailModule.getRemoveAdminAction();

  let newUserId = $state('');
</script>

<h3>Admins</h3>

{#await cobrandDetailModule.adminsAsync}
  <Tile>...getting admins</Tile>
{:then admins}
  <DataTable
    headers={[
      { key: 'user', value: 'User ID' },
      { key: 'id', empty: true },
    ]}
    rows={admins}
  >
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === 'id'}
        <ButtonSet>
          <Button kind="danger" onclick={() => removeAdminAction(cell.value)}>Remove</Button>
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
<TextInput labelText="User ID" bind:value={newUserId} />
<br />
<Button
  onclick={() => {
    addAdminAction(cobrandId, newUserId);
    newUserId = '';
  }}
>
  Add Admin
</Button>
