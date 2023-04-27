<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.ship.pr}
  {#if displayPattern === "full"}
    <td class="pr">
      {player.ship_stats.personal_rating.toFixed(0)}
    </td>
  {:else if displayPattern === "nopr"}
    <td class="pr" />
  {/if}
{/if}
