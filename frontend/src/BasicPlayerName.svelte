<script lang="ts">
import type { vo } from "wailsjs/go/models";
import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
import Const from "./Const";
export let player: vo.Player;

function clanURL(player: vo.Player): string {
  return (
    Const.BASE_NUMBERS_URL +
    "clan/" +
    player.player_info.clan.id +
    "," +
    player.player_info.clan.tag
  );
}

function playerURL(player: vo.Player): string {
  return (
    Const.BASE_NUMBERS_URL +
    "player/" +
    player.player_info.id +
    "," +
    player.player_info.name
  );
}
</script>

<td class="td-string omit">
  {#if player.player_info.id === 0}
    {player.player_info.name}
  {:else}
    <!-- svelte-ignore a11y-invalid-attribute -->
    {#if player.player_info.clan.id !== 0}
      <a
        class="td-link"
        id="clan-text"
        href="#"
        on:click="{() => BrowserOpenURL(clanURL(player))}"
      >
        [{player.player_info.clan.tag}]
      </a>
    {/if}
    <!-- svelte-ignore a11y-invalid-attribute -->
    <a
      class="td-link"
      href="#"
      on:click="{() => BrowserOpenURL(playerURL(player))}"
    >
      {player.player_info.name}
    </a>
  {/if}
</td>
