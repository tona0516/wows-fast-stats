<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";

export let player: vo.Player;
export let displayPattern: DisplayPattern;
</script>

{#if displayPattern === "full" || displayPattern === "nopr"}
  <td class="td-number">
    {player.ship_stats.battles}
  </td>
{/if}
