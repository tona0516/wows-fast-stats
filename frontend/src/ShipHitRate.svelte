<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
import type { StatsPattern } from "./StatsPattern";
import { values } from "./util";
import type { StatsCategory } from "./StatsCategory";

export let player: vo.Player;
export let displayPattern: DisplayPattern;
export let statsPattern: StatsPattern;
export let statsCatetory: StatsCategory;

const digit = Const.DIGITS["hit_rate"];

$: mainBattery = values(
  player,
  displayPattern,
  statsPattern,
  statsCatetory,
  "main_battery_hit_rate"
);
$: torpedoes = values(
  player,
  displayPattern,
  statsPattern,
  statsCatetory,
  "torpedoes_hit_rate"
);
</script>

{#if mainBattery !== undefined && torpedoes != undefined}
  <td class="td-multiple">
    {mainBattery.toFixed(digit)}% | {torpedoes.toFixed(digit)}%
  </td>
{/if}
