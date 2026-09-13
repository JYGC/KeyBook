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

  const backendClient = getBackendClient();
  const cobrandRepository = new CobrandRepository(backendClient);
  const cobrandAdminRepository = new CobrandAdminRepository(backendClient);
  const cobrandManagerRepo = new CobrandPropertyManagerRepository(backendClient);
  const cobrandService = new CobrandService(cobrandRepository, cobrandAdminRepository, cobrandManagerRepo);
  const cobrandListModule = new CobrandListModule(cobrandService);

  const gotoAddCobrand = () => {
    goto('/user/cobrands/add');
  };
</script>

<Button onclick={gotoAddCobrand}>Add Cobrand</Button>
<CobrandList {cobrandListModule} />
