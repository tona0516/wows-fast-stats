<script lang="ts">
  import type { DisplayPattern } from "../enums";

  export let shipColspan: number;
  export let overallColspan: number;
  export let displayPattern: DisplayPattern;

  $: allColspan = shipColspan + overallColspan;
</script>

{#if allColspan > 0}
  {#if displayPattern === "private"}
    <td class="no_data omit" colspan={allColspan}>PRIVATE</td>
  {:else if displayPattern === "nodata"}
    <td class="no_data omit" colspan={allColspan}>N/A</td>
  {/if}
{/if}

{#if shipColspan > 0}
  {#if displayPattern === "noshipstats"}
    <td class="no_data omit" colspan={shipColspan}>N/A</td>
  {/if}
{/if}
