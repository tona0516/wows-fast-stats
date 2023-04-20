<script lang="ts">
  import {
    GetConfig,
    GetTempArenaInfoHash,
    Load,
  } from "../wailsjs/go/main/App.js";
  import ConfigPage from "./ConfigPage.svelte";
  import StatsPage from "./MainPage.svelte";
  import type { vo } from "wailsjs/go/models.js";
  import Notification from "./Notification.svelte";
  import { WindowReloadApp } from "../wailsjs/runtime/runtime";

  type Page = "main" | "config";
  let currentPage: Page = "main";

  let loadState: LoadState;
  let latestHash: string;
  let teams: vo.Team[];
  let config: vo.UserConfig;

  let notification: Notification;

  setInterval(looper, 1000);

  async function looper() {
    try {
      config = await GetConfig();
      notification.removeToastWithKey("need_config");
    } catch (error) {
      notification.showToastWithKey(
        "未設定の状態のため開始できません。「設定」から入力してください。",
        "info",
        "need_config"
      );
      return;
    }

    if (loadState === "error" || loadState === "fetching") {
      notification.removeToastWithKey("not_in_battle");
      return;
    }

    let hash: string;
    try {
      hash = await GetTempArenaInfoHash();
    } catch (error) {
      loadState = "standby";
      notification.showToastWithKey(
        "戦闘中ではありません。開始時に自動的にリロードします。",
        "info",
        "not_in_battle"
      );
      return;
    }

    if (hash === latestHash) {
      return;
    }

    loadState = "fetching";
    try {
      const start = new Date().getTime();
      teams = await Load();
      latestHash = hash;
      loadState = "standby";
      const elapsed = (new Date().getTime() - start) / 1000;
      notification.showToast(`データ取得完了: ${elapsed}秒`, "success");
    } catch (error) {
      loadState = "error";
      notification.showToast(error, "error");
    }
  }

  function clickMain() {
    currentPage = "main";
  }

  function clickConfig() {
    currentPage = "config";
  }

  function clickReload() {
    WindowReloadApp();
  }
</script>

<main>
  <div style="font-size: {config?.font_size || 'medium'};">
    <Notification bind:this={notification} />

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
            <button
              type="button"
              class="btn btn-outline-secondary mx-1"
              title="ホーム"
              on:click={clickMain}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                fill="currentColor"
                class="bi bi-house"
                viewBox="0 0 16 16"
              >
                <path
                  d="M8.707 1.5a1 1 0 0 0-1.414 0L.646 8.146a.5.5 0 0 0 .708.708L2 8.207V13.5A1.5 1.5 0 0 0 3.5 15h9a1.5 1.5 0 0 0 1.5-1.5V8.207l.646.647a.5.5 0 0 0 .708-.708L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293L8.707 1.5ZM13 7.207V13.5a.5.5 0 0 1-.5.5h-9a.5.5 0 0 1-.5-.5V7.207l5-5 5 5Z"
                />
              </svg>
            </button>
            <button
              type="button"
              class="btn btn-outline-secondary mx-1"
              title="設定"
              on:click={clickConfig}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                fill="currentColor"
                class="bi bi-gear"
                viewBox="0 0 16 16"
              >
                <path
                  d="M8 4.754a3.246 3.246 0 1 0 0 6.492 3.246 3.246 0 0 0 0-6.492zM5.754 8a2.246 2.246 0 1 1 4.492 0 2.246 2.246 0 0 1-4.492 0z"
                />
                <path
                  d="M9.796 1.343c-.527-1.79-3.065-1.79-3.592 0l-.094.319a.873.873 0 0 1-1.255.52l-.292-.16c-1.64-.892-3.433.902-2.54 2.541l.159.292a.873.873 0 0 1-.52 1.255l-.319.094c-1.79.527-1.79 3.065 0 3.592l.319.094a.873.873 0 0 1 .52 1.255l-.16.292c-.892 1.64.901 3.434 2.541 2.54l.292-.159a.873.873 0 0 1 1.255.52l.094.319c.527 1.79 3.065 1.79 3.592 0l.094-.319a.873.873 0 0 1 1.255-.52l.292.16c1.64.893 3.434-.902 2.54-2.541l-.159-.292a.873.873 0 0 1 .52-1.255l.319-.094c1.79-.527 1.79-3.065 0-3.592l-.319-.094a.873.873 0 0 1-.52-1.255l.16-.292c.893-1.64-.902-3.433-2.541-2.54l-.292.159a.873.873 0 0 1-1.255-.52l-.094-.319zm-2.633.283c.246-.835 1.428-.835 1.674 0l.094.319a1.873 1.873 0 0 0 2.693 1.115l.291-.16c.764-.415 1.6.42 1.184 1.185l-.159.292a1.873 1.873 0 0 0 1.116 2.692l.318.094c.835.246.835 1.428 0 1.674l-.319.094a1.873 1.873 0 0 0-1.115 2.693l.16.291c.415.764-.42 1.6-1.185 1.184l-.291-.159a1.873 1.873 0 0 0-2.693 1.116l-.094.318c-.246.835-1.428.835-1.674 0l-.094-.319a1.873 1.873 0 0 0-2.692-1.115l-.292.16c-.764.415-1.6-.42-1.184-1.185l.159-.291A1.873 1.873 0 0 0 1.945 8.93l-.319-.094c-.835-.246-.835-1.428 0-1.674l.319-.094A1.873 1.873 0 0 0 3.06 4.377l-.16-.292c-.415-.764.42-1.6 1.185-1.184l.292.159a1.873 1.873 0 0 0 2.692-1.115l.094-.319z"
                />
              </svg>
            </button>
            <button
              type="button"
              class="btn btn-outline-success mx-1"
              title="リロード"
              on:click={clickReload}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                fill="currentColor"
                class="bi bi-arrow-clockwise"
                viewBox="0 0 16 16"
              >
                <path
                  fill-rule="evenodd"
                  d="M8 3a5 5 0 1 0 4.546 2.914.5.5 0 0 1 .908-.417A6 6 0 1 1 8 2v1z"
                />
                <path
                  d="M8 4.466V.534a.25.25 0 0 1 .41-.192l2.36 1.966c.12.1.12.284 0 .384L8.41 4.658A.25.25 0 0 1 8 4.466z"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </nav>

    {#if currentPage === "config"}
      <ConfigPage
        on:SuccessToast={(event) =>
          notification.showToast(event.detail.message, "success")}
        on:ErrorToast={(event) =>
          notification.showToast(event.detail.message, "error")}
      />
    {/if}

    {#if currentPage === "main"}
      <StatsPage bind:loadState bind:latestHash bind:teams bind:config />
    {/if}
  </div>
</main>
