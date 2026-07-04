<script lang="ts">
	import type { ICobrandDetailModule } from '$lib/modules/interfaces';
	import type { ICobrandAdminModel } from '$lib/models/cobrand-models';
	import { Button, ButtonSet, DataTable, TextInput, Tile } from 'carbon-components-svelte';

	let { cobrandId, cobrandDetailModule, currentUserId } = $props<{
		cobrandId: string;
		cobrandDetailModule: ICobrandDetailModule;
		currentUserId: string;
	}>();

	const addAdminAction = cobrandDetailModule.getAddAdminAction();
	const removeAdminAction = cobrandDetailModule.getRemoveAdminAction();

	let newUserId = $state('');
</script>

<h3>Admins</h3>

{#await cobrandDetailModule.adminsAsync}
	<Tile>...getting admins</Tile>
{:then admins}
	{@const ownAdmin = admins.find((a: ICobrandAdminModel) => a.user === currentUserId)}
	{#if ownAdmin && !ownAdmin.approved}
		<Tile>Pending approval — contact KeyBook to activate management access for this cobrand.</Tile>
		<br />
	{/if}
	<DataTable
		headers={[
			{ key: 'user', value: 'User ID' },
			{ key: 'approved', value: 'Approved' },
			{ key: 'id', empty: true }
		]}
		rows={admins}
	>
		<svelte:fragment slot="cell" let:cell>
			{#if cell.key === 'id'}
				<ButtonSet>
					<Button kind="danger" onclick={() => removeAdminAction(cell.value)}>Remove</Button>
				</ButtonSet>
			{:else if cell.key === 'approved'}
				{cell.value ? 'Yes' : 'Pending'}
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
