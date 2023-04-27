<script lang="ts">
  import { Cell } from "@smui/data-table";
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

<!-- ship kd rate -->
{#if config.displays.ship.kd_rate}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <Cell class="kd">
      {player.ship_stats.kd_rate.toFixed(1)}
    </Cell>
  {:else if displayPattern === "noshipstats"}
    <Cell class="kd" />
  {/if}
{/if}
