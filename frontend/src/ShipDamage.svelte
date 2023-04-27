<script lang="ts">
  import { Cell } from "@smui/data-table";
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.ship.damage}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <Cell class="damage">
      {player.ship_stats.avg_damage.toFixed(0)}
    </Cell>
  {:else if displayPattern === "noshipstats"}
    <Cell class="damage" />
  {/if}
{/if}
