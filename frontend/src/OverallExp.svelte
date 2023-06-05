<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";

export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["exp"];
</script>

{#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
  <td class="td-number">
    {player.overall_stats.exp.toFixed(digit)}
  </td>
{/if}
