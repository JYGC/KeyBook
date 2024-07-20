<script lang="ts">
  const { data } = $props();

  let resultList = $derived(data);
</script>

<p>
  {JSON.stringify(resultList)}
</p>
<a href="/importdata/csv1">import csv</a>