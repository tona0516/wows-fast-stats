<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";

export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["kd_rate"];
</script>

{#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
  <td class="td-number">
    {player.overall_stats.kd_rate.toFixed(digit)}
  </td>
{/if}
