<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  import { LogDebug } from "../wailsjs/runtime/runtime";
  export let config: vo.UserConfig;
  export let displayPattern: DisplayPattern;

  function countDisplays(config: vo.UserConfig): number {
    const shipCount = Object.values(config.displays.ship).filter(
      (it) => it === true
    ).length;
    const overallCount = Object.values(config.displays.overall).filter(
      (it) => it === true
    ).length;

    return shipCount + overallCount;
  }
</script>

{#if countDisplays(config) > 0}
  {#if displayPattern === "private"}
    <td class="no_data omit" colspan={countDisplays(config)}>PRIVATE</td>
  {:else if displayPattern === "nodata"}
    <td class="no_data omit" colspan={countDisplays(config)}>N/A</td>
  {/if}
{/if}
