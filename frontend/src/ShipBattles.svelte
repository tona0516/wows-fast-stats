<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import type { StatsPattern } from "./StatsPattern";
import { values } from "./util";
import type { StatsCategory } from "./StatsCategory";

export let player: vo.Player;
export let displayPattern: DisplayPattern;
export let statsPattern: StatsPattern;
export let statsCatetory: StatsCategory;

$: value = values(
  player,
  displayPattern,
  statsPattern,
  statsCatetory,
  "battles"
);
</script>

{#if value !== undefined}
  <td class="td-number">
    {value}
  </td>
{/if}
