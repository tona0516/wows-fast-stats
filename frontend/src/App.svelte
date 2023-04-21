<script lang="ts">
  import { WindowReloadApp } from "../wailsjs/runtime/runtime.js";
  import {
    GetConfig,
    GetTempArenaInfoHash,
    Load,
    SaveScreenshot,
  } from "../wailsjs/go/main/App.js";
  import type { vo } from "wailsjs/go/models.js";
  import Notification from "./Notification.svelte";
  import ConfigPage from "./ConfigPage.svelte";
  import MainPage from "./MainPage.svelte";
  import Navigation from "./Navigation.svelte";
  import domtoimage from "dom-to-image";
  import { LogDebug } from "../wailsjs/runtime/runtime.js";

  type Page = "main" | "config";
  let currentPage: Page = "main";

  let loadState: LoadState;
  let latestHash: string;
  let teams: vo.Team[];
  let config: vo.UserConfig;

  let notification: Notification;

  setInterval(looper, 1000);

  function onClickMenu(menu: NavigationMenu) {
    switch (menu) {
      case "main":
        currentPage = "main";
        break;
      case "config":
        currentPage = "config";
        break;
      case "reload":
        WindowReloadApp();
        break;
      default:
        break;
    }
  }

  async function sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

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

    if (config.save_screenshot) {
      try {
        const dataUrl = (await domtoimage.toPng(
          document.getElementById("mainpage")
        )) as string;
        const base64Data = dataUrl.split(",")[1];
        await SaveScreenshot(base64Data);
      } catch (error) {
        notification.showToast(error, "error");
      }
    }
  }
</script>

<main>
  <div style="font-size: {config?.font_size || 'medium'};">
    <Notification bind:this={notification} />

    <Navigation on:onClickMemu={(event) => onClickMenu(event.detail.menu)} />

    {#if currentPage === "config"}
      <ConfigPage
        on:SuccessToast={(event) =>
          notification.showToast(event.detail.message, "success")}
        on:ErrorToast={(event) =>
          notification.showToast(event.detail.message, "error")}
      />
    {/if}

    {#if currentPage === "main"}
      <div id="mainpage">
        <MainPage bind:loadState bind:latestHash bind:teams bind:config />
      </div>
    {/if}
  </div>
</main>
