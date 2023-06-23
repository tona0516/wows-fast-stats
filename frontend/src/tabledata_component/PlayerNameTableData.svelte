<script lang="ts">
import clone from "clone";
import { createEventDispatcher } from "svelte";
import type { vo } from "../../wailsjs/go/models";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import { SkillLevelConverter } from "../RankConverter";
import { clanURL, playerURL, values } from "../util";
import { StatsCategory } from "../enums";

export let player: vo.Player;
export let userConfig: vo.UserConfig;
export let statsPattern: string;
export let alertPlayers: vo.AlertPlayer[];

const dispatch = createEventDispatcher();

$: alertPlayer = alertPlayers.find(
  (it) => it.account_id === player.player_info.id
);
$: pr = values(player, statsPattern, StatsCategory.Ship, "pr");
$: color = SkillLevelConverter.fromPR(
  pr,
  userConfig.custom_color.skill
).toBgColorCode();
$: playerLabel = isBelongToClan()
  ? `[${player.player_info.clan.tag}] ${player.player_info.name}`
  : player.player_info.name;

function isBelongToClan(): boolean {
  return player.player_info.clan.id !== 0;
}
</script>

<td class="td-string omit" style="background-color: {color}">
  {#if player.player_info.id === 0}
    {playerLabel}
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
      {playerLabel}
    </a>

    <ul
      class="dropdown-menu"
      aria-labelledby="dropdownMenuLink"
      style="font-size: {userConfig.font_size};"
    >
      {#if isBelongToClan()}
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
