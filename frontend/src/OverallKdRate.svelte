<script lang="ts">
import type { vo } from "wailsjs/go/models";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;
</script>

{#if config.displays.overall.kd_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <td class="td-number">
      {player.overall_stats.kd_rate.toFixed(1)}
    </td>
  {/if}
{/if}
