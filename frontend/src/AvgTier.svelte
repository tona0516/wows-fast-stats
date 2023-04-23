<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

<!-- avg tier -->
{#if config.displays.player_avg_tier}
  {#if displayPattern === "noshipstats" || displayPattern === "full" || displayPattern === "nopr"}
    <td class="avg-tier">
      {player.player_stats.avg_tier.toFixed(1)}
    </td>
  {/if}
{/if}
