<script lang="ts">
  import { Cell } from "@smui/data-table";
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.ship.pr}
  {#if displayPattern === "full"}
    <Cell class="pr">
      {player.ship_stats.personal_rating.toFixed(0)}
    </Cell>
  {:else if displayPattern === "noshipstats" || displayPattern === "nopr"}
    <Cell class="pr" />
  {/if}
{/if}
