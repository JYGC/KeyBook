<script lang="ts">
  import { goto } from '$app/navigation';
  import { getBackendClient } from '$lib/api/backend-client';
  import ItemList from '$lib/components/item/ItemList.svelte';
  import { ItemRepository } from '$lib/repositories/item/item-repository';
  import { EntryDeviceRepository } from '$lib/repositories/item/entry-device-repository';
  import { ItemService } from '$lib/services/item/item-service';
  import { ItemListModule } from '$lib/modules/item/item-list-module.svelte';
  import { Button } from 'carbon-components-svelte';

  const pb = getBackendClient();
  const itemRepo = new ItemRepository(pb);
  const entryDeviceRepo = new EntryDeviceRepository(pb);
  const itemService = new ItemService(itemRepo, entryDeviceRepo);
  const itemListModule = new ItemListModule(itemService);

  const gotoAddItem = () => {
    goto('/user/items/add');
  };
</script>

<Button onclick={gotoAddItem}>Add Item</Button>
<ItemList {itemListModule} />
