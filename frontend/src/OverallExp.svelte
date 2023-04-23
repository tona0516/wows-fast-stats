<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.player_exp}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td class="exp">
      {player.player_stats.exp.toFixed(0)}
    </td>
  {:else if displayPattern === "noshipstats"}
    <td class="exp" />
  {/if}
{/if}
