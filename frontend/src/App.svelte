<script lang="ts">
import {
  UserConfig,
  TempArenaInfoHash,
  Battle,
  SaveScreenshot,
} from "../wailsjs/go/main/App.js";
import Notification from "./Notification.svelte";
import ConfigPage from "./PageConfig.svelte";
import MainPage from "./PageMain.svelte";
import domtoimage from "dom-to-image";
import { LogDebug, WindowReloadApp } from "../wailsjs/runtime/runtime.js";
import AppInfo from "./PageAppInfo.svelte";

import "bootstrap-icons/font/bootstrap-icons.css";
import PageHelp from "./PageHelp.svelte";
import type { vo } from "wailsjs/go/models.js";

type Page = "main" | "config" | "help" | "appinfo";
type Func = "reload" | "screenshot";
type NavigationMenu = Page | Func;
type ScreenshotType = "auto" | "manual";

let currentPage: Page = "main";

let loadState: LoadState;
let latestHash: string;
let battle: vo.Battle;
let config: vo.UserConfig;

let notification: Notification;

let timer = setInterval(looper, 1000);

function onClickMenu(menu: NavigationMenu) {
  switch (menu) {
    case "main":
      currentPage = "main";
      break;
    case "config":
      currentPage = "config";
      break;
    case "help":
      currentPage = "help";
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
        notification.showToast("スクリーンショットを保存しました。", "success");
      }
    })
    .catch((_) => {});
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

  clearInterval(timer);
  loadState = "fetching";
  try {
    const start = new Date().getTime();
    battle = await Battle();
    latestHash = hash;
    loadState = "standby";

    const elapsed = (new Date().getTime() - start) / 1000;
    notification.showToast(`データ取得完了: ${elapsed}秒`, "success");
    timer = setInterval(looper, 1000);
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
        <button
          class="navbar-toggler"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#navbarNavAltMarkup"
          aria-controls="navbarNavAltMarkup"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
          <div class="navbar-nav">
            <button
              type="button"
              class="btn btn-sm btn-outline-secondary m-1 {currentPage ===
                'main' && 'active'}"
              title="ホーム"
              on:click="{() => onClickMenu('main')}"
            >
              <i class="bi bi-house"></i>
              ホーム
            </button>
            <button
              type="button"
              class="btn btn-sm btn-outline-secondary m-1 {currentPage ===
                'config' && 'active'}"
              title="設定"
              on:click="{() => onClickMenu('config')}"
            >
              <i class="bi bi-gear"></i>
              設定
            </button>
            <button
              type="button"
              class="btn btn-sm btn-outline-secondary m-1 {currentPage ===
                'help' && 'active'}"
              title="設定"
              on:click="{() => onClickMenu('help')}"
            >
              <i class="bi bi-question-circle"></i>
              ヘルプ
            </button>
            <button
              type="button"
              class="btn btn-sm btn-outline-secondary m-1 {currentPage ===
                'appinfo' && 'active'}"
              title="アプリ情報"
              on:click="{() => onClickMenu('appinfo')}"
            >
              <i class="bi bi-info-circle"></i>
              アプリ情報
            </button>
            {#if currentPage == "main"}
              <button
                type="button"
                class="btn btn-sm btn-outline-success m-1"
                title="リロード"
                on:click="{() => onClickMenu('reload')}"
              >
                <i class="bi bi-arrow-clockwise"></i>
                リロード
              </button>

              <button
                type="button"
                class="btn btn-sm btn-outline-success m-1"
                title="スクリーンショット"
                disabled="{battle === undefined || loadState === 'fetching'}"
                on:click="{() => onClickMenu('screenshot')}"
              >
                <i class="bi bi-camera"></i>
                スクリーンショット
              </button>
            {/if}
          </div>
        </div>
      </div>
    </nav>

    {#if currentPage === "main"}
      <div id="mainpage">
        <MainPage
          bind:loadState="{loadState}"
          bind:latestHash="{latestHash}"
          bind:battle="{battle}"
          bind:config="{config}"
        />
      </div>
    {/if}

    {#if currentPage === "config"}
      <ConfigPage
        on:SuccessToast="{(event) =>
          notification.showToast(event.detail.message, 'success')}"
        on:ErrorToast="{(event) =>
          notification.showToast(event.detail.message, 'error')}"
      />
    {/if}

    {#if currentPage === "help"}
      <PageHelp />
    {/if}

    {#if currentPage === "appinfo"}
      <AppInfo />
    {/if}

    <Notification bind:this="{notification}" />
  </div>
</main>
