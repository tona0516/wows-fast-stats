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

let digit = Const.DIGITS["survived_rate"];

$: win = values(
  player,
  displayPattern,
  statsPattern,
  statsCatetory,
  "win_survived_rate"
);
$: lose = values(
  player,
  displayPattern,
  statsPattern,
  statsCatetory,
  "lose_survived_rate"
);
</script>

{#if win !== undefined && lose !== undefined}
  <td class="td-multiple">
    {win.toFixed(digit)}% | {lose.toFixed(digit)}%
  </td>
{/if}
