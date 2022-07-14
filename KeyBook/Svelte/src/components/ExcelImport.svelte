<svelte:options tag="excel-import-svelte" />
<script lang="ts">
</script>
<main>
    <form action="/DataImport/Excel" enctype="multipart/form-data" method="POST">
        <label for="spreadsheet">Spreatsheet</label>
        <input type="file" name="postedfile" id="postedfile">
        <button>Upload spreadsheet</button>
    </form>
</main>