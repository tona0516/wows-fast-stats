<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.ship.win_survived_rate}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td class="survived-rate">
      {player.ship_stats.win_survived_rate.toFixed(1)}
    </td>
  {/if}
{/if}
