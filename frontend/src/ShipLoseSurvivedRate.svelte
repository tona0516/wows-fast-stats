<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.ship.lose_survived_rate}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td class="td-number">
      {player.ship_stats.lose_survived_rate.toFixed(1)}
    </td>
  {/if}
{/if}
