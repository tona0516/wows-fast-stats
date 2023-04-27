<script lang="ts">
  import { Cell } from "@smui/data-table";
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.ship.battles}
  {#if displayPattern === "full" || displayPattern === "nopr"}
    <Cell class="battles">
      {player.ship_stats.battles}
    </Cell>
  {:else if displayPattern === "noshipstats"}
    <Cell class="battles" />
  {/if}
{/if}
