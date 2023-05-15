<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["survived_rate"];
</script>

{#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
  <td class="td-multiple">
    {player.overall_stats.win_survived_rate.toFixed(digit)}% | {player.overall_stats.lose_survived_rate.toFixed(
      digit
    )}%
  </td>
{/if}
