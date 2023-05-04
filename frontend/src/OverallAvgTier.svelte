<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;
</script>

<!-- avg tier -->
{#if config.displays.overall.avg_tier}
  {#if displayPattern === "noshipstats" || displayPattern === "full" || displayPattern === "nopr"}
    <td class="td-number">
      {player.overall_stats.avg_tier.toFixed(1)}
    </td>
  {/if}
{/if}
