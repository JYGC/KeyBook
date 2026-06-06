<script lang="ts">
  import { goto } from '$app/navigation';
  import { getBackendClient } from '$lib/api/backend-client';
  import CobrandList from '$lib/components/cobrand/CobrandList.svelte';
  import { CobrandRepository } from '$lib/repositories/cobrand/cobrand-repository';
  import { CobrandAdminRepository } from '$lib/repositories/cobrand/cobrand-admin-repository';
  import { CobrandPropertyManagerRepository } from '$lib/repositories/cobrand/cobrand-property-manager-repository';
  import { CobrandService } from '$lib/services/cobrand/cobrand-service';
  import { CobrandListModule } from '$lib/modules/cobrand/cobrand-list-module.svelte';
  import { Button } from 'carbon-components-svelte';

  const pb = getBackendClient();
  const cobrandRepo = new CobrandRepository(pb);
  const cobrandAdminRepo = new CobrandAdminRepository(pb);
  const cobrandManagerRepo = new CobrandPropertyManagerRepository(pb);
  const cobrandService = new CobrandService(cobrandRepo, cobrandAdminRepo, cobrandManagerRepo);
  const cobrandListModule = new CobrandListModule(cobrandService);

  const gotoAddCobrand = () => {
    goto('/user/cobrands/add');
  };
</script>

<Button onclick={gotoAddCobrand}>Add Cobrand</Button>
<CobrandList {cobrandListModule} />
