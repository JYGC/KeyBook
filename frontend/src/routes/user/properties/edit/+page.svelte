<script lang="ts">
	import { goto } from "$app/navigation";
	import { BackendClient } from "$lib/api/backend-client";
	import PropertyEditor from "$lib/components/property/PropertyEditor.svelte";
	import ConfirmButtonAndDialog from "$lib/components/shared/ConfirmButtonAndDialog.svelte";
	import { getPropertyContext } from "$lib/contexts/property-context.svelte";
	import type { IEditPropertyModel } from "$lib/models/property-models";
	import { Button, Tile } from "carbon-components-svelte";

  const propertyContext = getPropertyContext();

  const backendClient = new BackendClient();

  let propertyAsync = $derived.by<Promise<IEditPropertyModel | null>>(async () => {
    try {
      return await backendClient.pb.collection("properties").getOne<IEditPropertyModel>(
        propertyContext.selectedPropertyId,
        { fields: "id,address,owners,managers" },
      );
    } catch (ex) {
      alert(ex);
      return null;
    }
  });

  const gotoPropertyList = () => {
    goto("/user/properties/list");
  }

  const savePropertiesActionAsync = async (changedProperty: IEditPropertyModel) => {
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

  const deletePersonActionAsync = async (property: IEditPropertyModel) => {
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
    >
      {#snippet deleteButton(deleteActionButtonClick: () => null)}
        <ConfirmButtonAndDialog
        submitAction={deleteActionButtonClick}
        buttonText="Delete"
        bodyMessage="Are you sure you want to delete this property?"
        />
      {/snippet}
    </PropertyEditor>
  {/if}
{:catch error}
  {error}
{/await}