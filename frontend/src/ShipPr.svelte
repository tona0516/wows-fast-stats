<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
import { RankConverter } from "./RankConverter";
export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["pr"];

$: color = RankConverter.fromPR(player.ship_stats.pr).toTextColorCode();
</script>

{#if displayPattern === "full"}
  <td class="td-number" style="color: {color}">
    {player.ship_stats.pr.toFixed(digit)}
  </td>
{:else if displayPattern === "nopr"}
  <td class="td-number"></td>
{/if}
