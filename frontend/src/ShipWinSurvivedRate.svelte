<script lang="ts">
  import { Cell } from "@smui/data-table";
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.ship.win_survived_rate}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <Cell class="survived-rate">
      {player.ship_stats.win_survived_rate.toFixed(1)}
    </Cell>
  {:else if displayPattern === "noshipstats"}
    <Cell class="survived-rate" />
  {/if}
{/if}
