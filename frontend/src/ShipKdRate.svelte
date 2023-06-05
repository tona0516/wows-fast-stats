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

const digit = Const.DIGITS["kd_rate"];

$: value = values(
  player,
  displayPattern,
  statsPattern,
  statsCatetory,
  "kd_rate"
);
</script>

{#if value !== undefined}
  <td class="td-number">
    {value.toFixed(digit)}
  </td>
{/if}
