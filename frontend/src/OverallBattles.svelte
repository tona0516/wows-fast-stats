<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";

export let player: vo.Player;
export let displayPattern: DisplayPattern;
</script>

{#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
  <td class="td-number">
    {player.overall_stats.battles}
  </td>
{/if}
