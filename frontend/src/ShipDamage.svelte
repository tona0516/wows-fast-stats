<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
import { RankConverter } from "./RankConverter";
import type { StatsPattern } from "./StatsPattern";
import { values } from "./util";
import type { StatsCategory } from "./StatsCategory";

export let player: vo.Player;
export let displayPattern: DisplayPattern;
export let statsPattern: StatsPattern;
export let statsCatetory: StatsCategory;

const digit = Const.DIGITS["damage"];

$: color = RankConverter.fromDamage(
  value,
  player.ship_info.avg_damage
).toTextColorCode();
$: value = values(
  player,
  displayPattern,
  statsPattern,
  statsCatetory,
  "damage"
);
</script>

{#if value !== undefined}
  <td class="td-number" style="color: {color}">
    {value.toFixed(digit)}
  </td>
{/if}
