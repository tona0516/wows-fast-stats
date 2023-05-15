<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
import RankConverter from "./RankConverter";
export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["win_rate"];

$: color = RankConverter.fromWinRate(
  player.overall_stats.win_rate
).toTextColorCode();
</script>

{#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
  <td class="td-number" style="color: {color}">
    {player.overall_stats.win_rate.toFixed(digit)}%
  </td>
{/if}
