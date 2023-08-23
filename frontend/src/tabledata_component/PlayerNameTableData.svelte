<script lang="ts">
  import clone from "clone";
  import { SkillLevelConverter } from "src/RankConverter";
  import { StatsCategory } from "src/enums";
  import { storedUserConfig } from "src/stores";
  import { values, clanURL, playerURL } from "src/util";
  import { createEventDispatcher } from "svelte";
  import type { domain } from "wailsjs/go/models";
  import { BrowserOpenURL } from "wailsjs/runtime/runtime";

  export let player: domain.Player;
  export let userConfig: domain.UserConfig;
  export let statsPattern: string;
  export let alertPlayers: domain.AlertPlayer[];

  const dispatch = createEventDispatcher();

  $: alertPlayer = alertPlayers.find(
    (it) => it.account_id === player.player_info.id
  );
  $: bgcolor = bgColor(player, userConfig, statsPattern);
  $: playerLabel = isBelongToClan(player)
    ? `[${player.player_info.clan.tag}] ${player.player_info.name}`
    : player.player_info.name;

  function isBelongToClan(player: domain.Player): boolean {
    return player.player_info.clan.id !== 0;
  }

  function bgColor(
    player: domain.Player,
    userConfig: domain.UserConfig,
    statsPattern: string
  ): string {
    let statsCategory: StatsCategory;
    if (userConfig.custom_color.player_name === "ship") {
      statsCategory = StatsCategory.Ship;
    }
    if (userConfig.custom_color.player_name === "overall") {
      statsCategory = StatsCategory.Overall;
    }

    if (!statsCategory) {
      return "";
    }

    const pr = values(player, statsPattern, statsCategory, "pr");
    return SkillLevelConverter.fromPR(
      pr,
      userConfig.custom_color.skill
    ).toBgColorCode();
  }
</script>

<td class="td-string omit" style="background-color: {bgcolor}">
  {#if player.player_info.id === 0}
    {playerLabel}
  {:else}
    {#if alertPlayer}
      <i class="bi {alertPlayer.pattern}" />
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
      style="font-size: {$storedUserConfig.font_size};"
    >
      {#if isBelongToClan(player)}
        <!-- svelte-ignore a11y-invalid-attribute -->
        <li>
          <a
            class="dropdown-item"
            href="#"
            on:click={() => BrowserOpenURL(clanURL(player))}
            >クラン詳細(WoWS Stats & Numbers)</a
          >
        </li>
      {/if}
      <!-- svelte-ignore a11y-invalid-attribute -->
      <li>
        <a
          class="dropdown-item"
          href="#"
          on:click={() => BrowserOpenURL(playerURL(player))}
          >プレイヤー詳細(WoWS Stats & Numbers)</a
        >
      </li>
      <!-- svelte-ignore a11y-invalid-attribute -->
      <li>
        {#if alertPlayer}
          <a
            class="dropdown-item"
            href="#"
            on:click={() =>
              dispatch("RemoveAlertPlayer", { target: clone(alertPlayer) })}
            >プレイヤーリストから削除する</a
          >
        {:else}
          <a
            class="dropdown-item"
            href="#"
            on:click={() => {
              dispatch("UpdateAlertPlayer", {
                target: {
                  account_id: player.player_info.id,
                  name: player.player_info.name,
                  pattern: "bi-check-circle-fill",
                  message: "",
                },
              });
            }}>プレイヤーリストへ追加する</a
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
