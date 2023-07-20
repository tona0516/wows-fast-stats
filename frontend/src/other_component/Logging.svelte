<script lang="ts">
  export let logs: string[];
</script>

<div class="center">
  <div class="logging overflow-scroll bg-dark p-2 mx-2">
    {#each logs as log}
      {log}
      <br />
    {/each}
  </div>
</div>
