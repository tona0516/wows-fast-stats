<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.overall.win_survived_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <td class="td-number">
      {player.player_stats.win_survived_rate.toFixed(1)}
    </td>
  {/if}
{/if}
