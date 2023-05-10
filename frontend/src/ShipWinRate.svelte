<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["win_rate"];
</script>

{#if config.displays.ship.win_rate}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td class="td-number">
      {player.ship_stats.win_rate.toFixed(digit)}%
    </td>
  {/if}
{/if}
