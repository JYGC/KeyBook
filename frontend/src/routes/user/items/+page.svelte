<script lang="ts">
  import { goto } from '$app/navigation';
  import { getBackendClient } from '$lib/api/backend-client';
  import ItemList from '$lib/components/item/ItemList.svelte';
  import { ItemRepository } from '$lib/repositories/item/item-repository';
  import { EntryDeviceRepository } from '$lib/repositories/item/entry-device-repository';
  import { ItemService } from '$lib/services/item/item-service';
  import { ItemListModule } from '$lib/modules/item/item-list-module.svelte';
  import { Button } from 'carbon-components-svelte';

  const backendClient = getBackendClient();
  const itemRepository = new ItemRepository(backendClient);
  const entryDeviceRepository = new EntryDeviceRepository(backendClient);
  const itemService = new ItemService(itemRepository, entryDeviceRepository);
  const itemListModule = new ItemListModule(itemService);

  const gotoAddItem = () => {
    goto('/user/items/add');
  };
</script>

<Button onclick={gotoAddItem}>Add Item</Button>
<ItemList {itemListModule} />
