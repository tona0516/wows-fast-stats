<script lang="ts">
  import { Cell } from "@smui/data-table";
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.overall.battles}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <Cell class="battles">
      {player.player_stats.battles}
    </Cell>
  {/if}
{/if}
