<script lang="ts">
import type { vo } from "wailsjs/go/models";
import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
import Const from "./Const";
import { RankConverter } from "./RankConverter";
import { createEventDispatcher } from "svelte";
import clone from "clone";
import { storedAlertPlayers, storedUserConfig } from "./stores";
import { get } from "svelte/store";

export let player: vo.Player;

let alertPlayers = get(storedAlertPlayers);
storedAlertPlayers.subscribe((it) => (alertPlayers = it));

let userConfig = get(storedUserConfig);
storedUserConfig.subscribe((it) => (userConfig = it));

const dispatch = createEventDispatcher();

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

$: alertPlayer = alertPlayers.find(
  (it) => it.account_id === player.player_info.id
);
$: color = RankConverter.fromPR(player.ship_stats.pr).toBgColorCode();
</script>

<td class="td-string omit" style="background-color: {color}">
  {#if player.player_info.id === 0}
    {player.player_info.name}
  {:else}
    {#if alertPlayer}
      <i class="bi {alertPlayer.pattern}"></i>
    {/if}

    <!-- svelte-ignore a11y-invalid-attribute -->
    <a
      class="td-link dropdown-toggle"
      href="#"
      id="dropdownMenuLink"
      data-bs-toggle="dropdown"
    >
      {#if player.player_info.clan.id !== 0}
        [{player.player_info.clan.tag}] {player.player_info.name}
      {:else}
        {player.player_info.name}
      {/if}
    </a>

    <ul
      class="dropdown-menu"
      aria-labelledby="dropdownMenuLink"
      style="font-size: {userConfig.font_size};"
    >
      {#if player.player_info.clan.id !== 0}
        <!-- svelte-ignore a11y-invalid-attribute -->
        <li>
          <a
            class="dropdown-item"
            href="#"
            on:click="{() => BrowserOpenURL(clanURL(player))}"
            >クラン詳細(WoWS Stats & Numbers)</a
          >
        </li>
      {/if}
      <!-- svelte-ignore a11y-invalid-attribute -->
      <li>
        <a
          class="dropdown-item"
          href="#"
          on:click="{() => BrowserOpenURL(playerURL(player))}"
          >プレイヤー詳細(WoWS Stats & Numbers)</a
        >
      </li>
      <!-- svelte-ignore a11y-invalid-attribute -->
      <li>
        {#if alertPlayer}
          <a
            class="dropdown-item"
            href="#"
            on:click="{() =>
              dispatch('RemoveAlertPlayer', { target: clone(alertPlayer) })}"
            >プレイヤーリストから削除する</a
          >
        {:else}
          <a
            class="dropdown-item"
            href="#"
            on:click="{() => {
              dispatch('UpdateAlertPlayer', {
                target: {
                  account_id: player.player_info.id,
                  name: player.player_info.name,
                  pattern: 'bi-check-circle-fill',
                  message: '',
                },
              });
            }}">プレイヤーリストへ追加する</a
          >
        {/if}
      </li>
      {#if alertPlayer}
        <li>
          <div class="dropdown-item">メモ: {alertPlayer.message}</div>
        </li>
      {/if}
    </ul>
  {/if}
</td>
