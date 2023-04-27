<script lang="ts">
  import { Cell } from "@smui/data-table";
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.overall.kd_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <Cell class="kd">
      {player.player_stats.kd_rate.toFixed(1)}
    </Cell>
  {/if}
{/if}
