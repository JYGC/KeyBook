<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import PropertyEditor from "$lib/components/property/PropertyEditor.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import type { IEditPropertyDto } from "$lib/dtos/property-dtos";
	import { Button, Tile } from "carbon-components-svelte";

  const propertyContext = getPropertyContext();

  const backendClient = new BackendClient();

  let propertyAsync = $derived.by<Promise<IEditPropertyDto | null>>(async () => {
    try {
      return await backendClient.pb.collection("properties").getOne<IEditPropertyDto>(
        propertyContext.selectedPropertyId,
        { fields: "id,address,owners,managers" },
      );
    } catch (ex) {
      alert(ex);
      return null;
    }
  });

  const gotoPropertyList = () => {
    goto("/properties/list");
  }

  const savePropertiesActionAsync = async (changedProperty: IEditPropertyDto) => {
    try {
      await backendClient.pb.collection("properties").update(changedProperty.id, {
        address: changedProperty.address,
        owners: changedProperty.owners,
        managers: changedProperty.managers,
      });
      gotoPropertyList();
    } catch (ex) {
      alert(ex);
    }
  };

  const deletePersonActionAsync = async (property: IEditPropertyDto) => {
    try {
      await backendClient.pb.collection("properties").delete(property.id);
      gotoPropertyList();
    } catch (ex) {
      alert(ex);
    }
  };
</script>

<Button onclick={gotoPropertyList}>Back</Button>

{#await propertyAsync}
  <Tile>...getting property details</Tile>
{:then property}
  {#if property === null}
    {history.back()}
  {:else}
    <PropertyEditor
      property={property}
      isAdd={false}
      savePropertyAction={savePropertiesActionAsync}
      deletePropertyAction={deletePersonActionAsync}
    />
  {/if}
{:catch error}
  {error}
{/await}