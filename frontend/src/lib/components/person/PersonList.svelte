<script lang="ts">
	import type { IPersonListItemDto } from "$lib/dtos/person-dtos";
	import { type IBackendClient } from "../../interfaces";
  
  let {
    backendClient,
    propertyId,
  } = $props<{
    backendClient: IBackendClient,
    propertyId: string,
  }>();

  let personListAsync = $derived.by<Promise<IPersonListItemDto>>(async () => {
    try {
      const response = await backendClient.pb.collection("persons").getList(1, 50, {
        filter: `property.id = "${propertyId}"`,
        fields: "id,type,name",
      });
      return response.items;
    } catch (ex) {
      alert(ex);
      return [];
    }
  });
</script>

{#await personListAsync}
  ...getting persons
{:then personList}
  {JSON.stringify(personList)}
{:catch error}
  {error}
{/await}