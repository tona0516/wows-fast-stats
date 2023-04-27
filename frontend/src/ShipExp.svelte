<script lang="ts">
  import { Cell } from "@smui/data-table";
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.ship.exp}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <Cell class="exp">
      {player.ship_stats.exp.toFixed(0)}
    </Cell>
  {:else if displayPattern === "noshipstats"}
    <Cell class="exp" />
  {/if}
{/if}
