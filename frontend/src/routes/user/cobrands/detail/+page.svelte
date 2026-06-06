<script lang="ts">
  import { getBackendClient } from '$lib/api/backend-client';
  import CobrandEditor from '$lib/components/cobrand/CobrandEditor.svelte';
  import CobrandAdminList from '$lib/components/cobrand/CobrandAdminList.svelte';
  import CobrandPropertyManagerList from '$lib/components/cobrand/CobrandPropertyManagerList.svelte';
  import ConfirmButtonAndDialog from '$lib/components/shared/ConfirmButtonAndDialog.svelte';
  import { CobrandRepository } from '$lib/repositories/cobrand/cobrand-repository';
  import { CobrandAdminRepository } from '$lib/repositories/cobrand/cobrand-admin-repository';
  import { CobrandPropertyManagerRepository } from '$lib/repositories/cobrand/cobrand-property-manager-repository';
  import { CobrandService } from '$lib/services/cobrand/cobrand-service';
  import { CobrandDetailModule } from '$lib/modules/cobrand/cobrand-detail-module.svelte';
  import { Button } from 'carbon-components-svelte';

  let { data } = $props();
  const cobrandId: string = data.cobrandId;

  const goBack = () => {
    history.back();
  };

  const pb = getBackendClient();
  const cobrandRepo = new CobrandRepository(pb);
  const cobrandAdminRepo = new CobrandAdminRepository(pb);
  const cobrandManagerRepo = new CobrandPropertyManagerRepository(pb);
  const cobrandService = new CobrandService(cobrandRepo, cobrandAdminRepo, cobrandManagerRepo);
  const cobrandDetailModule = new CobrandDetailModule(cobrandService, cobrandId, goBack);
</script>

<Button onclick={goBack}>Back</Button>

<CobrandEditor cobrandEditorModule={cobrandDetailModule}>
  {#snippet deleteButton(deleteActionButtonClick: () => void)}
    <ConfirmButtonAndDialog
      submitAction={deleteActionButtonClick}
      buttonText="Delete Cobrand"
      bodyMessage="Are you sure you want to delete this cobrand?"
    />
  {/snippet}
</CobrandEditor>

<br />

<CobrandAdminList {cobrandId} {cobrandDetailModule} />

<br />

<CobrandPropertyManagerList {cobrandId} {cobrandDetailModule} />
