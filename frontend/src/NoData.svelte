<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  import { LogDebug } from "../wailsjs/runtime/runtime";
  import { Cell } from "@smui/data-table";
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
    <Cell class="no_data omit" colspan={countDisplays(config)}>PRIVATE</Cell>
  {:else if displayPattern === "nodata"}
    <Cell class="no_data omit" colspan={countDisplays(config)}>N/A</Cell>
  {/if}
{/if}
