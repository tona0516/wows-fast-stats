<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
import { RankConverter } from "./RankConverter";
import { values } from "./util";
import type { StatsPattern } from "./StatsPattern";
import type { StatsCategory } from "./StatsCategory";

export let player: vo.Player;
export let displayPattern: DisplayPattern;
export let statsPattern: StatsPattern;
export let statsCatetory: StatsCategory;

const digit = Const.DIGITS["pr"];

$: color = RankConverter.fromPR(value).toTextColorCode();
$: value = values(player, displayPattern, statsPattern, statsCatetory, "pr");
</script>

{#if value !== undefined}
  <td class="td-number" style="color: {color}">
    {#if value !== -1}
      {value.toFixed(digit)}
    {/if}
  </td>
{/if}
