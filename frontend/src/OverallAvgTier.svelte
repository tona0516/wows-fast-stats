<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["avg_tier"];
</script>

<!-- avg tier -->
{#if config.displays.overall.avg_tier}
  {#if displayPattern === "noshipstats" || displayPattern === "full" || displayPattern === "nopr"}
    <td class="td-number">
      {player.overall_stats.avg_tier.toFixed(digit)}
    </td>
  {/if}
{/if}
