<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

<td class="td-string omit">
  <!-- svelte-ignore a11y-invalid-attribute -->
  <a href="#" on:click={() => BrowserOpenURL(player.player_info.stats_url)}>
    {#if player.player_info.clan}
      [{player.player_info.clan}]{player.player_info.name}
    {:else}
      {player.player_info.name}
    {/if}
  </a>
</td>
