<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
import TextColor from "./TextColor";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["win_rate"];
</script>

{#if config.displays.overall.win_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <td
      class="td-number"
      style="color: {TextColor.winRate(player.overall_stats.win_rate)}"
    >
      {player.overall_stats.win_rate.toFixed(digit)}%
    </td>
  {/if}
{/if}
