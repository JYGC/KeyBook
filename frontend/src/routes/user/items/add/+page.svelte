<script lang="ts">
  import { getBackendClient } from '$lib/api/backend-client';
  import ItemEditor from '$lib/components/item/ItemEditor.svelte';
  import { ItemRepository } from '$lib/repositories/item/item-repository';
  import { EntryDeviceRepository } from '$lib/repositories/item/entry-device-repository';
  import { ItemService } from '$lib/services/item/item-service';
  import { ItemAddModule } from '$lib/modules/item/item-add-module.svelte';
  import { Button } from 'carbon-components-svelte';

  const goBack = () => {
    history.back();
  };

  const pb = getBackendClient();
  const itemRepo = new ItemRepository(pb);
  const entryDeviceRepo = new EntryDeviceRepository(pb);
  const itemService = new ItemService(itemRepo, entryDeviceRepo);
  const itemAddModule = new ItemAddModule(itemService, goBack);
</script>

<Button onclick={goBack}>Back</Button>
<ItemEditor itemEditorModule={itemAddModule} />
