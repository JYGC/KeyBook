<script lang="ts">
  import { getBackendClient } from '$lib/api/backend-client';
  import CobrandEditor from '$lib/components/cobrand/CobrandEditor.svelte';
  import { CobrandRepository } from '$lib/repositories/cobrand/cobrand-repository';
  import { CobrandAdminRepository } from '$lib/repositories/cobrand/cobrand-admin-repository';
  import { CobrandPropertyManagerRepository } from '$lib/repositories/cobrand/cobrand-property-manager-repository';
  import { CobrandService } from '$lib/services/cobrand/cobrand-service';
  import { CobrandAddModule } from '$lib/modules/cobrand/cobrand-add-module.svelte';
  import { Button } from 'carbon-components-svelte';

  const goBack = () => {
    history.back();
  };

  const pb = getBackendClient();
  const currentUserId = pb.authStore.record?.id ?? '';
  const cobrandRepo = new CobrandRepository(pb);
  const cobrandAdminRepo = new CobrandAdminRepository(pb);
  const cobrandManagerRepo = new CobrandPropertyManagerRepository(pb);
  const cobrandService = new CobrandService(cobrandRepo, cobrandAdminRepo, cobrandManagerRepo);
  const cobrandAddModule = new CobrandAddModule(cobrandService, currentUserId, goBack);
</script>

<Button onclick={goBack}>Back</Button>
<CobrandEditor cobrandEditorModule={cobrandAddModule} />
