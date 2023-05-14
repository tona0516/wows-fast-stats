<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
import TextColor from "./TextColor";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["pr"];
</script>

{#if config.displays.ship.pr}
  {#if displayPattern === "full"}
    <td
      class="td-number"
      style="color: {TextColor.prText(player.ship_stats.pr)}"
    >
      {player.ship_stats.pr.toFixed(digit)}
    </td>
  {:else if displayPattern === "nopr"}
    <td class="td-number"></td>
  {/if}
{/if}
