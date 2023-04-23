<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  import { LogDebug } from "../wailsjs/runtime/runtime";
  export let config: vo.UserConfig;
  export let displayPattern: DisplayPattern;

  function countDisplays(config: vo.UserConfig): number {
    const excludedColumns = ["player_name", "ship_info"];

    return Object.entries(config.displays).filter(
      ([k, v]) => !excludedColumns.includes(k) && v === true
    ).length;
  }
</script>

{#if countDisplays(config) > 0}
  {#if displayPattern === "private"}
    <td class="no_data omit" colspan={countDisplays(config)}>PRIVATE</td>
  {:else if displayPattern === "nodata"}
    <td class="no_data omit" colspan={countDisplays(config)}>N/A</td>
  {/if}
{/if}
