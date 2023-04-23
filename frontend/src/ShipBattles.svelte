<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.ship_battles}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <td class="battles">
      {player.ship_stats.battles}
    </td>
  {:else if displayPattern === "noshipstats"}
    <td class="battles" />
  {/if}
{/if}
