<script lang="ts">
  import { Cell } from "@smui/data-table";
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

<!-- avg tier -->
{#if config.displays.overall.avg_tier}
  {#if displayPattern === "noshipstats" || displayPattern === "full" || displayPattern === "nopr"}
    <Cell class="avg-tier">
      {player.player_stats.avg_tier.toFixed(1)}
    </Cell>
  {/if}
{/if}
