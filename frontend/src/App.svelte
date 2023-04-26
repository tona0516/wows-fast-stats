<script lang="ts">
  import {
    UserConfig,
    TempArenaInfoHash,
    Battle,
    SaveScreenshot,
  } from "../wailsjs/go/main/App.js";
  import type { vo } from "wailsjs/go/models.js";
  import Notification from "./Notification.svelte";
  import ConfigPage from "./ConfigPage.svelte";
  import MainPage from "./MainPage.svelte";
  import domtoimage from "dom-to-image";
  import { LogDebug, WindowReloadApp } from "../wailsjs/runtime/runtime.js";
  import HomeIcon from "./HomeIcon.svelte";
  import ConfigIcon from "./ConfigIcon.svelte";
  import ReloadIcon from "./ReloadIcon.svelte";
  import CameraIcon from "./CameraIcon.svelte";
  import iconApp from "./assets/images/appicon.png";
  import InfoIcon from "./InfoIcon.svelte";
  import AppInfo from "./AppInfo.svelte";

  type NavigationMenu = "main" | "config" | "appinfo" | "reload" | "screenshot";
  type ScreenshotType = "auto" | "manual";
  type Page = "main" | "config" | "appinfo";

  let currentPage: Page = "main";

  let loadState: LoadState;
  let latestHash: string;
  let battle: vo.Battle;
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
      case "appinfo":
        currentPage = "appinfo";
        break;
      case "reload":
        WindowReloadApp();
        break;
      case "screenshot":
        saveScreenshot("manual");
      default:
        break;
    }
  }

  function saveScreenshot(type: ScreenshotType) {
    domtoimage
      .toPng(document.getElementById("mainpage"))
      .then((dataUrl) => {
        const date = battle.meta.date.replaceAll(":", "-").replaceAll(" ", "-");
        const ownShip = battle.meta.own_ship.replaceAll(" ", "-");
        const filename = `${date}_${ownShip}_${battle.meta.arena}_${battle.meta.type}.png`;
        const base64Data = dataUrl.split(",")[1];
        if (type === "auto") {
          return SaveScreenshot(filename, base64Data, false);
        }
        if (type === "manual") {
          return SaveScreenshot(filename, base64Data, true);
        }
      })
      .then(() => {
        if (type === "manual") {
          notification.showToast(
            "スクリーンショットを保存しました。",
            "success"
          );
        }
      })
      .catch((error) => {
        notification.showToast(error, "error");
      });
  }

  async function looper() {
    try {
      config = await UserConfig();
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
      hash = await TempArenaInfoHash();
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
      battle = await Battle();
      latestHash = hash;
      loadState = "standby";

      const elapsed = (new Date().getTime() - start) / 1000;
      notification.showToast(`データ取得完了: ${elapsed}秒`, "success");
    } catch (error) {
      loadState = "error";
      notification.showToast(error, "error");
    }

    if (config.save_screenshot) {
      await saveScreenshot("auto");
    }
  }
</script>

<main>
  <div style="font-size: {config?.font_size || 'medium'};">
    <nav class="navbar navbar-expand-sm navbar-light bg-light">
      <div class="container-fluid">
        <img class="me-2" src={iconApp} alt="" width="36px" height="36px" />
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
              class="btn btn-outline-secondary m-1 {currentPage === 'main' &&
                'active'}"
              title="ホーム"
              on:click={() => onClickMenu("main")}
            >
              <HomeIcon />
            </button>
            <button
              type="button"
              class="btn btn-outline-secondary m-1 {currentPage === 'config' &&
                'active'}"
              title="設定"
              on:click={() => onClickMenu("config")}
            >
              <ConfigIcon />
            </button>
            <button
              type="button"
              class="btn btn-outline-secondary m-1 {currentPage === 'appinfo' &&
                'active'}"
              title="アプリ情報"
              on:click={() => onClickMenu("appinfo")}
            >
              <InfoIcon />
            </button>
            {#if currentPage == "main"}
              <button
                type="button"
                class="btn btn-outline-success m-1"
                title="リロード"
                on:click={() => onClickMenu("reload")}
              >
                <ReloadIcon />
              </button>

              <button
                type="button"
                class="btn btn-outline-success m-1"
                title="スクリーンショット"
                disabled={battle === undefined || loadState === "fetching"}
                on:click={() => onClickMenu("screenshot")}
              >
                <CameraIcon />
              </button>
            {/if}
          </div>
        </div>
      </div>
    </nav>

    {#if currentPage === "main"}
      <div id="mainpage">
        <MainPage bind:loadState bind:latestHash bind:battle bind:config />
      </div>
    {/if}

    {#if currentPage === "config"}
      <ConfigPage
        on:SuccessToast={(event) =>
          notification.showToast(event.detail.message, "success")}
        on:ErrorToast={(event) =>
          notification.showToast(event.detail.message, "error")}
      />
    {/if}

    {#if currentPage === "appinfo"}
      <AppInfo />
    {/if}

    <Notification bind:this={notification} />
  </div>
</main>
