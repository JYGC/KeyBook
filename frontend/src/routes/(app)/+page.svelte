<script lang="ts">
  const { data } = $props();

  let resultList = $derived(data);
</script>

<p>
  {JSON.stringify(resultList)}
</p>
<a href="/devices/csvimport">import csv</a>