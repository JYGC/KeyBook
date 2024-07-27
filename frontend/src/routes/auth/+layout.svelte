<script lang="ts">
  const { children } = $props();
</script>

<div class="center">
  {@render children()}
</div>

<style>
  .center {
    position: fixed;
    inset: 0px;
    width: 24rem;
    height: 30rem;
    max-width: 100vw;
    max-height: 100dvh;
    margin: auto;
  }
</style>