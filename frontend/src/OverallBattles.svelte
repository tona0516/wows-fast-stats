<script lang="ts">
import type { vo } from "wailsjs/go/models";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;
</script>

{#if config.displays.overall.battles}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <td class="td-number">
      {player.overall_stats.battles}
    </td>
  {/if}
{/if}
