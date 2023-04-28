<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.overall.damage}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <td class="td-number">
      {player.player_stats.damage.toFixed(0)}
    </td>
  {/if}
{/if}
