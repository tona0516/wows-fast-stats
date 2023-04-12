<script lang="ts">
  import SvelteTable from "svelte-table";
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

  let friendRows = [];
  let enemyRows = [];
  let columns = [
    {
      title: "クラン",
      value: (v) => v.player_player_info.clan,
      class: "text-left",
    },
    {
      title: "プレイヤー",
      value: (v) => v.player_player_info.name,
      class: "text-left omit",
    },
    {
      title: "CP",
      value: (v) =>
        isValidStatsValue(v, "ship")
          ? v.player_ship_stats.combat_power.toFixed(0)
          : "",
      class: "text-right",
    },
    {
      title: "PR",
      value: (v) =>
        isValidStatsValue(v, "ship")
          ? v.player_ship_stats.personal_rating.toFixed(0)
          : "",
      class: "text-right",
    },
    {
      title: "艦",
      value: (v) => v.player_ship_info.name,
      class: "text-left omit",
    },
    {
      title: "T",
      value: (v) =>
        v.player_ship_info.tier === 11 ? "★" : v.player_ship_info.tier,
      class: "text-right",
    },
    {
      title: "Dmg",
      value: (v) =>
        isValidStatsValue(v, "ship")
          ? v.player_ship_stats.avg_damage.toFixed(0)
          : "",
      class: "text-right",
    },
    {
      title: "勝率",
      value: (v) =>
        isValidStatsValue(v, "ship")
          ? v.player_ship_stats.win_rate.toFixed(1)
          : "",
      class: "text-right",
    },
    {
      title: "Exp",
      value: (v) =>
        isValidStatsValue(v, "ship")
          ? v.player_ship_stats.avg_exp.toFixed(0)
          : "",
      class: "text-right",
    },
    {
      title: "戦闘数",
      value: (v) =>
        isValidStatsValue(v, "ship") ? v.player_ship_stats.battles : "",
      class: "text-right",
    },
    {
      title: "Dmg",
      value: (v) =>
        isValidStatsValue(v, "player")
          ? v.player_player_stats.avg_damage.toFixed(0)
          : "",
      class: "text-right",
    },
    {
      title: "勝率",
      value: (v) =>
        isValidStatsValue(v, "player")
          ? v.player_player_stats.win_rate.toFixed(1)
          : "",
      class: "text-right",
    },
    {
      title: "Exp",
      value: (v) =>
        isValidStatsValue(v, "player")
          ? v.player_player_stats.avg_exp.toFixed(0)
          : "",
      class: "text-right",
    },
    {
      title: "戦闘数",
      value: (v) =>
        isValidStatsValue(v, "player") ? v.player_player_stats.battles : "",
      class: "text-right",
    },
    {
      title: "平均T",
      value: (v) =>
        isValidStatsValue(v, "player")
          ? v.player_player_stats.avg_tier.toFixed(1)
          : "",
      class: "text-right",
    },
  ];

  setInterval(looper, 1000);

  async function looper() {
    try {
      await GetConfig();
      removeSettingPromotionIfNeeded()
    } catch (error) {
      showSettingPromotionIfNeeded()
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
      const stats = await Load();
      friendRows = stats["friends"];
      enemyRows = stats["enemies"];
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
      duration: 3000,
      placement: "bottom-right",
      type: "success",
      theme: "dark",
    });
  }

  function showErrorToast(message: string) {
    toasts.add({
      description: message,
      duration: 3000,
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
      description: "未設定の状態のため開始できません。「設定」から入力してください。",
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
    <div class="mt-3">
      {#if state === "fetching"}
        <div class="d-flex justify-content-center">
          <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
      {/if}

      {#if latestHash !== ""}
        <div class="padding">
          <SvelteTable
            {columns}
            rows={friendRows}
            classNameTable="table table-sm table-dark table-striped"
          />
          <SvelteTable
            {columns}
            rows={enemyRows}
            classNameTable="table table-sm table-dark table-striped"
          />
        </div>
      {/if}
    </div>
  {/if}

  <ToastContainer let:data>
    <FlatToast {data} />
  </ToastContainer>
</main>

<style>
  :global(.text-right) {
    text-align: right;
  }
  :global(.text-left) {
    text-align: left;
  }
  :global(.omit) {
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100px;
    white-space: nowrap;
  }
  :global(.padding) {
    padding: 1em;
  }
  :global(.horizontal) {
    display: flex;
  }
  :global(.form-style) {
    width: 50%;
    margin: auto;
  }
</style>
