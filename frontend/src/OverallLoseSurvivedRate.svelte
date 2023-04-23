<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.player_lose_survived_rate}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td class="survived-rate">
      {player.player_stats.lose_survived_rate.toFixed(1)}
    </td>
  {:else if displayPattern === "noshipstats"}
    <td class="survived-rate" />
  {/if}
{/if}
