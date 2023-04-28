<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.ship.exp}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td class="td-number">
      {player.ship_stats.exp.toFixed(0)}
    </td>
  {/if}
{/if}
