<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["hit_rate"];
</script>

{#if config.displays.ship.hit_rate}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td class="td-multiple">
      {player.ship_stats.main_battery_hit_rate.toFixed(digit)}% | {player.ship_stats.torpedoes_hit_rate.toFixed(
        digit
      )}%
    </td>
  {/if}
{/if}
