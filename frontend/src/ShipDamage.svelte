<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
import TextColor from "./TextColor";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["damage"];
</script>

{#if config.displays.ship.damage}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td
      class="td-number"
      style="color: {TextColor.shipDamage(
        player.ship_stats.damage,
        player.ship_info.avg_damage
      )}"
    >
      {player.ship_stats.damage.toFixed(digit)}
    </td>
  {/if}
{/if}
