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

  const addHouseholdMemberAction = propertyDetailModule.getAddHouseholdMemberAction();
  const removeHouseholdMemberAction = propertyDetailModule.getRemoveHouseholdMemberAction();

  let newPersonId = $state('');
</script>

<h3>Household Members</h3>

{#await propertyDetailModule.householdMembersAsync}
  <Tile>...getting household members</Tile>
{:then householdMembers}
  <DataTable
    headers={[
      { key: 'person', value: 'Person ID' },
      { key: 'id', empty: true },
    ]}
    rows={householdMembers}
  >
    <svelte:fragment slot="cell" let:cell>
      {#if cell.key === 'id'}
        <ButtonSet>
          <Button kind="danger" onclick={() => removeHouseholdMemberAction(cell.value)}>Remove</Button>
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
    addHouseholdMemberAction(newPersonId, propertyId);
    newPersonId = '';
  }}
>
  Add Household Member
</Button>
