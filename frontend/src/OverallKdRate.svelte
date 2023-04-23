<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.player_kd_rate}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td class="kd">
      {player.player_stats.kd_rate.toFixed(1)}
    </td>
  {:else if displayPattern === "noshipstats"}
    <td class="kd" />
  {/if}
{/if}
