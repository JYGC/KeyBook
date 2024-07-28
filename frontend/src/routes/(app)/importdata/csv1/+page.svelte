<script lang="ts">
	import { CsvFileToObjectConverter } from "$lib/modules/csv-file-to-object-converter.svelte";
	import { UploadCsvApi } from "$lib/api/upload-csv-api";
	import { BackendClient } from "$lib/api/backend-client.svelte";
	import { goto } from "$app/navigation";
	import { Button } from "carbon-components-svelte";

  const acceptedExtensions = ['.csv']
  
  const csvFileToObjectConverter = new CsvFileToObjectConverter();
  const authManager = new BackendClient();
  const uploadCsvApi = new UploadCsvApi(csvFileToObjectConverter, authManager);

  const uploadFile = () => {
    uploadCsvApi.callApi();
    csvFileToObjectConverter.input = null;
  };

  let disableUploadButtonAsync = $derived.by(async () => await csvFileToObjectConverter.outputAsync === null);

  const gotoPropertyList = () => {
    goto("/properties/list");
  };
</script>

<Button onclick={gotoPropertyList}>Back</Button>

<input
  type="file"
  name="fileToUpload"
  id="file"
  accept={acceptedExtensions.join(',')}
  bind:files={csvFileToObjectConverter.input}
  required
/>
{#await disableUploadButtonAsync then value} 
  <Button
    disabled={value}
    onclick={uploadFile}
  >Upload</Button>
{/await}
<p>
  {#await csvFileToObjectConverter.outputAsync}
    ...awaiting
  {:then output}
    {#if output !== null}
      {JSON.stringify(output)}
    {/if}
  {:catch error}
    {error}
  {/await}
</p>