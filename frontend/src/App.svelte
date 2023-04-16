<script lang="ts">
  import {
    Debug,
    GetConfig,
    GetTempArenaInfoHash,
    Load,
  } from "../wailsjs/go/main/App.js";
  import ConfigPage from "./ConfigPage.svelte";
  import StatsPage from "./MainPage.svelte";
  import type { vo } from "wailsjs/go/models.js";
  import Notification from "./Notification.svelte";

  type Page = "main" | "config";
  let currentPage: Page = "main";

  let loadState: LoadState;
  let latestHash: string;
  let teams: vo.Team[];
  let notification: Notification;

  let config: vo.UserConfig;

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
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
              class="nav-link"
              href="#"
              data-bs-toggle="collapse"
              data-bs-target=".navbar-collapse.show"
              on:click={clickMain}>ホーム</a
            >
            <!-- svelte-ignore a11y-invalid-attribute -->
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
      <ConfigPage
        on:SuccessToast={(event) =>
          notification.showToast(event.detail.message, "success")}
        on:ErrorToast={(event) =>
          notification.showToast(event.detail.message, "error")}
      />
    {/if}

    {#if currentPage === "main"}
      <StatsPage bind:loadState bind:latestHash bind:teams />
    {/if}
  </div>
</main>
