<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.player_win_rate}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td class="win">
      {player.player_stats.win_rate.toFixed(1)}
    </td>
  {:else if displayPattern === "noshipstats"}
    <td class="win" />
  {/if}
{/if}
