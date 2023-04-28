<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.ship.pr}
  {#if displayPattern === "full"}
    <td class="td-number">
      {player.ship_stats.pr.toFixed(0)}
    </td>
  {:else if displayPattern === "nopr"}
    <td class="td-number" />
  {/if}
{/if}
