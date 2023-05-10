<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["kd_rate"];
</script>

<!-- ship kd rate -->
{#if config.displays.ship.kd_rate}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td class="td-number">
      {player.ship_stats.kd_rate.toFixed(digit)}
    </td>
  {/if}
{/if}
