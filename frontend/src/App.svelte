<script lang="ts">
  import { toasts, ToastContainer, FlatToast } from "svelte-toasts";
  import { fade } from "svelte/transition";
  import {
    Debug,
    GetConfig,
    GetTempArenaInfoHash,
    Load,
  } from "../wailsjs/go/main/App.js";
  import type { ToastProps } from "svelte-toasts/types/common.js";
  import ConfigPage from "./ConfigPage.svelte";
  import StatsPage from "./StatsPage.svelte";
  import type { vo } from "wailsjs/go/models.js";

  type Page = "main" | "config";
  let currentPage: Page = "main";

  let notInBattleToast: ToastProps;
  let settingPromotionToast: ToastProps;

  let state;
  let latestHash;
  let teams;

  let config: vo.UserConfig;

  setInterval(looper, 1000);

  async function looper() {
    try {
      config = await GetConfig();
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
</script>

<main>
  <div style="font-size: {config?.font_size || 'medium'};">
    <ToastContainer let:data>
      <FlatToast {data} />
    </ToastContainer>

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
      <ConfigPage
        on:SuccessToast={(event) => showSuccessToast(event.detail.message)}
        on:ErrorToast={(event) => showErrorToast(event.detail.message)}
      />
    {/if}

    {#if currentPage === "main"}
      <StatsPage bind:state bind:latestHash bind:teams />
    {/if}
  </div>
</main>

<style>
  :global(.aligner) {
    display: flex;
    align-items: center;
  }
</style>
