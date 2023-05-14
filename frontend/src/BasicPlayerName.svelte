<script lang="ts">
import type { vo } from "wailsjs/go/models";
import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
export let player: vo.Player;
</script>

<td class="td-string omit">
  {#if player.player_info.id === 0}
    {player.player_info.name}
  {:else}
    <!-- svelte-ignore a11y-invalid-attribute -->
    <a
      class="td-link"
      href="#"
      on:click="{() => BrowserOpenURL(player.player_info.stats_url)}"
    >
      {#if player.player_info.clan}
        [{player.player_info.clan}]{player.player_info.name}
      {:else}
        {player.player_info.name}
      {/if}
    </a>
  {/if}
</td>
