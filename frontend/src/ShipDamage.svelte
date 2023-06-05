<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
import { RankConverter } from "./RankConverter";

export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["damage"];

$: color = RankConverter.fromDamage(
  player.ship_stats.damage,
  player.ship_info.avg_damage
).toTextColorCode();
</script>

{#if displayPattern === "full" || displayPattern === "nopr"}
  <td class="td-number" style="color: {color}">
    {player.ship_stats.damage.toFixed(digit)}
  </td>
{/if}
