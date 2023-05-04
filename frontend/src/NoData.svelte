<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
export let config: vo.UserConfig;
export let displayPattern: DisplayPattern;

const shipCounts = Object.values(config.displays.ship).filter(
  (it) => it === true
).length;

const overallCounts = Object.values(config.displays.overall).filter(
  (it) => it === true
).length;
</script>

{#if shipCounts + overallCounts > 0}
  {#if displayPattern === "private"}
    <td class="no_data omit" colspan="{shipCounts + overallCounts}">PRIVATE</td>
  {:else if displayPattern === "nodata"}
    <td class="no_data omit" colspan="{shipCounts + overallCounts}">N/A</td>
  {:else if displayPattern === "noshipstats"}
    <td class="no_data omit" colspan="{shipCounts}">N/A</td>
  {/if}
{/if}
