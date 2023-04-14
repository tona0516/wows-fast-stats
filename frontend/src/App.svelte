<script lang="ts">
  import iconCV from "./assets/images/icon-cv.png";
  import iconBB from "./assets/images/icon-bb.png";
  import iconCL from "./assets/images/icon-cl.png";
  import iconDD from "./assets/images/icon-dd.png";
  import iconSS from "./assets/images/icon-ss.png";
  import iconNone from "./assets/images/icon-none.png";
  import { BrowserOpenURL } from "../wailsjs/runtime";

  import { toasts, ToastContainer, FlatToast } from "svelte-toasts";
  import {
    ApplyConfig,
    Debug,
    GetConfig,
    GetTempArenaInfoHash,
    Load,
    SelectDirectory,
  } from "../wailsjs/go/main/App.js";
  import type { ToastProps } from "svelte-toasts/types/common.js";

  type Page = "main" | "config";
  let currentPage: Page = "main";

  type State = "standby" | "fetching" | "error";
  let state: State = "standby";

  type StatsType = "ship" | "player";

  let notInBattleToast: ToastProps;
  let settingPromotionToast: ToastProps;
  let latestHash = "";

  let installPath = "";
  let appid = "";

  let teams = [];
  setInterval(looper, 1000);

  async function looper() {
    try {
      await GetConfig();
      removeSettingPromotionIfNeeded();
    } catch (error) {
      showSettingPromotionIfNeeded();
      return;
    }

    if (state === "error" || state === "fetching") {
      removeNotInBattleToastIfNeeded();
      return;
    }

    let hash: string;
    try {
      hash = await GetTempArenaInfoHash();
    } catch (error) {
      state = "standby";
      showNotInBattleToastIfNeeded();
      return;
    }

    if (hash === latestHash) {
      return;
    }

    state = "fetching";
    try {
      const start = new Date().getTime();
      teams = await Load();
      console.log(teams);
      latestHash = hash;
      state = "standby";
      const elapsed = (new Date().getTime() - start) / 1000;
      showSuccessToast(`データ取得完了: ${elapsed}秒`);
    } catch (error) {
      state = "error";
      showErrorToast(error);
    }
  }

  function clickMain() {
    currentPage = "main";
  }

  function clickConfig() {
    currentPage = "config";
    GetConfig()
      .then((config) => {
        installPath = config.install_path;
        appid = config.appid;
      })
      .catch((error) => {
        installPath = "";
        appid = "";
      });
  }

  function clickApply() {
    ApplyConfig(installPath, appid)
      .then((_) => {
        showSuccessToast("更新しました。");
      })
      .catch((error) => {
        showErrorToast(error);
      });
  }

  function selectDirectory() {
    SelectDirectory().then((result) => {
      if (!result) return;
      installPath = result;
    });
  }

  function showSuccessToast(message: string) {
    toasts.add({
      description: message,
      duration: 5000,
      placement: "bottom-right",
      type: "success",
      theme: "dark",
    });
  }

  function showErrorToast(message: string) {
    toasts.add({
      description: message,
      duration: 5000,
      placement: "bottom-right",
      type: "error",
      theme: "dark",
    });
  }

  function showNotInBattleToastIfNeeded() {
    if (notInBattleToast) {
      return;
    }

    notInBattleToast = toasts.add({
      description: "戦闘中ではありません。開始時に自動的にリロードします。",
      duration: 0,
      placement: "bottom-right",
      type: "info",
      theme: "dark",
    });
  }

  function removeNotInBattleToastIfNeeded() {
    if (!notInBattleToast) {
      return;
    }

    notInBattleToast.remove();
    notInBattleToast = undefined;
  }

  function showSettingPromotionIfNeeded() {
    if (settingPromotionToast) {
      return;
    }

    settingPromotionToast = toasts.add({
      description:
        "未設定の状態のため開始できません。「設定」から入力してください。",
      duration: 0,
      placement: "bottom-right",
      type: "info",
      theme: "dark",
    });
  }

  function removeSettingPromotionIfNeeded() {
    if (!settingPromotionToast) {
      return;
    }

    settingPromotionToast.remove();
    settingPromotionToast = undefined;
  }

  function isValidStatsValue(player: any, statsType: StatsType) {
    let battles: number;
    switch (statsType) {
      case "ship":
        battles = player.player_ship_stats.battles;
        break;
      case "player":
        battles = player.player_player_stats.battles;
        break;
      default:
        battles = 0;
        break;
    }
    return !player.player_player_info.is_hidden && battles > 0;
  }

  function backgroundClassForTr(personalRating: number): string {
    console.log(personalRating);
    switch (true) {
      case personalRating == 0:
        return "";
      case personalRating < 750:
        return "bad";
      case personalRating < 1100:
        return "below-average";
      case personalRating < 1350:
        return "average";
      case personalRating < 1550:
        return "good";
      case personalRating < 1750:
        return "very-good";
      case personalRating < 2100:
        return "great";
      case personalRating < 2450:
        return "unicum";
      case personalRating >= 2450:
        return "super-unicum";
      default:
        return "";
    }
  }

  function numberForDisplay(value: number): string {
    if (value === 11) return "★";

    const decimal = [10, 9, 5, 4, 1];
    const romanNumeral = ["X", "IX", "V", "IV", "I"];

    let romanized = "";

    for (var i = 0; i < decimal.length; i++) {
      while (decimal[i] <= value) {
        romanized += romanNumeral[i];
        value -= decimal[i];
      }
    }
    return romanized;
  }

  function shipIconForDisplay(shipType: string): string {
    switch (shipType) {
      case "AirCarrier":
        return iconCV;
      case "Battleship":
        return iconBB;
      case "Cruiser":
        return iconCL;
      case "Destroyer":
        return iconDD;
      case "Submarine":
        return iconSS;
      default:
        return iconNone;
    }
  }
</script>

<main>
  <nav class="navbar navbar-expand-lg navbar-light bg-light">
    <div class="container-fluid">
      <span class="navbar-brand">wows-fast-stats</span>
      <button
        class="navbar-toggler"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarNavAltMarkup"
        aria-controls="navbarNavAltMarkup"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span class="navbar-toggler-icon" />
      </button>
      <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
        <div class="navbar-nav">
          <a
            class="nav-link"
            href="#"
            data-bs-toggle="collapse"
            data-bs-target=".navbar-collapse.show"
            on:click={clickMain}>ホーム</a
          >
          <a
            class="nav-link"
            href="#"
            data-bs-toggle="collapse"
            data-bs-target=".navbar-collapse.show"
            on:click={clickConfig}>設定</a
          >
        </div>
      </div>
    </div>
  </nav>

  {#if currentPage === "config"}
    <div class="mt-3">
      <form>
        <div class="mb-3 form-style">
          <label for="install-path" class="form-label"
            >World of Warshipsインストールフォルダ</label
          >
          <div class="horizontal">
            <input
              type="text"
              class="form-control"
              id="install-path"
              bind:value={installPath}
            />
            <button
              type="button"
              class="btn btn-secondary"
              on:click={selectDirectory}>選択</button
            >
          </div>
        </div>
        <div class="mb-3 form-style">
          <label for="appid" class="form-label">AppID</label>
          <input
            type="text"
            class="form-control"
            id="appid"
            bind:value={appid}
          />
        </div>
        <button type="button" class="btn btn-primary" on:click={clickApply}
          >適用</button
        >
      </form>
    </div>
  {/if}

  {#if currentPage === "main"}
    <div>
      {#if state === "fetching"}
        <div class="d-flex justify-content-center m-3">
          <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
      {/if}

      {#if latestHash !== ""}
        <div class="mt-1 mx-3">
          {#each teams as team}
            <table class="table table-sm">
              <thead>
                <tr>
                  <th>プレイヤー</th>
                  <th>PR</th>
                  <th>艦</th>
                  <th>Dmg(艦)</th>
                  <th>勝率(艦)</th>
                  <th>K/D(艦)</th>
                  <th>戦闘数(艦)</th>
                  <th>Dmg</th>
                  <th>勝率</th>
                  <th>K/D</th>
                  <th>戦闘数</th>
                  <th>平均T</th>
                </tr>
              </thead>
              <tbody>
                {#each team as player}
                  <tr
                    class={backgroundClassForTr(
                      player.player_ship_stats.personal_rating
                    )}
                  >
                    <td class="name omit">
                      <a
                        href="#"
                        on:click={BrowserOpenURL(
                          player.player_player_info.stats_url
                        )}
                      >
                        {#if player.player_player_info.clan}
                          [{player.player_player_info.clan}]{player
                            .player_player_info.name}
                        {:else}
                          {player.player_player_info.name}
                        {/if}
                      </a>
                    </td>
                    {#if isValidStatsValue(player, "ship")}
                      <td class="pr"
                        >{player.player_ship_stats.personal_rating.toFixed(
                          0
                        )}</td
                      >
                    {:else}
                      <td class="pr" />
                    {/if}

                    <td class="name omit">
                      <a
                        href="#"
                        on:click={BrowserOpenURL(
                          player.player_ship_info.stats_url
                        )}
                      >
                        <div class="aligner">
                          <img
                            alt=""
                            width="24px"
                            height="24px"
                            src={shipIconForDisplay(
                              player.player_ship_info.type
                            )}
                          />
                          <div class="omit">
                            {numberForDisplay(player.player_ship_info.tier)}
                            {player.player_ship_info.name}
                          </div>
                        </div>
                      </a>
                    </td>
                    {#if isValidStatsValue(player, "ship")}
                      <td class="damage"
                        >{player.player_ship_stats.avg_damage.toFixed(0)}</td
                      >
                    {:else}
                      <td class="damage" />
                    {/if}

                    {#if isValidStatsValue(player, "ship")}
                      <td class="win"
                        >{player.player_ship_stats.win_rate.toFixed(1)}</td
                      >
                    {:else}
                      <td class="win" />
                    {/if}

                    {#if isValidStatsValue(player, "ship")}
                      <td class="kd"
                        >{player.player_ship_stats.kd_rate.toFixed(1)}</td
                      >
                    {:else}
                      <td class="kd" />
                    {/if}

                    {#if isValidStatsValue(player, "ship")}
                      <td class="battles">{player.player_ship_stats.battles}</td
                      >
                    {:else}
                      <td class="battles" />
                    {/if}

                    {#if isValidStatsValue(player, "player")}
                      <td class="damage"
                        >{player.player_player_stats.avg_damage.toFixed(0)}</td
                      >
                    {:else}
                      <td class="damage" />
                    {/if}
                    {#if isValidStatsValue(player, "player")}
                      <td class="win"
                        >{player.player_player_stats.win_rate.toFixed(1)}</td
                      >
                    {:else}
                      <td class="win" />
                    {/if}
                    {#if isValidStatsValue(player, "player")}
                      <td class="kd"
                        >{player.player_player_stats.kd_rate.toFixed(1)}</td
                      >
                    {:else}
                      <td class="kd" />
                    {/if}
                    {#if isValidStatsValue(player, "player")}
                      <td class="battles"
                        >{player.player_player_stats.battles}</td
                      >
                    {:else}
                      <td class="battles" />
                    {/if}
                    {#if isValidStatsValue(player, "player")}
                      <td class="avg-tier"
                        >{player.player_player_stats.avg_tier.toFixed(1)}</td
                      >
                    {:else}
                      <td class="avg-tier" />
                    {/if}
                  </tr>
                {/each}
              </tbody>
            </table>
          {/each}
        </div>
      {/if}
    </div>
  {/if}

  <ToastContainer let:data>
    <FlatToast {data} />
  </ToastContainer>
</main>

<style>
  :global(.horizontal) {
    display: flex;
  }
  :global(.form-style) {
    width: 50%;
    margin: auto;
  }

  :global(.aligner) {
    display: flex;
    align-items: center;
  }
</style>
