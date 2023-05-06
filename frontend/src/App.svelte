<script lang="ts">
import {
  UserConfig,
  Battle,
  ManualScreenshot,
  AutoScreenshot,
  ExcludePlayerIDs,
  Prepare,
  IsFinishedPrepare,
  Ready,
} from "../wailsjs/go/main/App.js";
import Notification from "./Notification.svelte";
import ConfigPage from "./PageConfig.svelte";
import MainPage from "./PageMain.svelte";
import { toPng } from "html-to-image";
import {
  EventsOn,
  LogInfo,
  WindowReloadApp,
} from "../wailsjs/runtime/runtime.js";
import AppInfo from "./PageAppInfo.svelte";

import "bootstrap-icons/font/bootstrap-icons.css";
import PageHelp from "./PageHelp.svelte";
import type { vo } from "wailsjs/go/models.js";
import { Average, type AverageFactor } from "./Average.js";

type Page = "main" | "config" | "help" | "appinfo";
type Func = "reload" | "screenshot";
type NavigationMenu = Page | Func;

let currentPage: Page = "main";

let battle: vo.Battle;
let config: vo.UserConfig;
let averageFactors: AverageFactor;
let excludePlayerIDs: number[];

let notification: Notification;

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
      manualScreenshot();
    default:
      break;
  }
}

async function getScreenshotBase64(): Promise<[string, string]> {
  const dataUrl = await toPng(document.getElementById("mainpage"));
  const date = battle.meta.date.replaceAll(":", "-").replaceAll(" ", "-");
  const ownShip = battle.meta.own_ship.replaceAll(" ", "-");
  const filename = `${date}_${ownShip}_${battle.meta.arena}_${battle.meta.type}.png`;
  const base64Data = dataUrl.split(",")[1];
  return [filename, base64Data];
}

async function manualScreenshot() {
  const [filename, data] = await getScreenshotBase64();
  try {
    await ManualScreenshot(filename, data);
    notification.showToast("スクリーンショットを保存しました。", "success");
  } catch (error) {
    notification.showToast("スクリーンショットの保存に失敗しました。", "error");
  }
}

async function autoScreenshot() {
  const [filename, data] = await getScreenshotBase64();
  await AutoScreenshot(filename, data);
}

EventsOn("BATTLE_START", async () => {
  try {
    notification.removeToastWithKey("not_in_battle");
    notification.showToastWithKey("戦闘データの取得中...", "info", "battle");

    const start = new Date().getTime();

    battle = await Battle();
    excludePlayerIDs = await ExcludePlayerIDs();
    const average = new Average(battle);
    averageFactors = average.calc(excludePlayerIDs);

    const elapsed = (new Date().getTime() - start) / 1000;
    notification.showToast(`データ取得完了: ${elapsed}秒`, "success");

    if (config.save_screenshot) {
      autoScreenshot();
    }
  } catch (error) {
    notification.showToastWithKey(error, "error", "error");
  } finally {
    notification.removeToastWithKey("battle");
  }
});

EventsOn("BATTLE_END", () => {
  notification.showToastWithKey(
    "戦闘中ではありません。開始時に自動的にリロードします。",
    "info",
    "not_in_battle"
  );
});

EventsOn("BATTLE_ERROR", (error) => {
  notification.showToastWithKey(error, "error", "error");
});

async function main() {
  try {
    config = await UserConfig();
  } catch (error) {
    notification.showToastWithKey(
      "未設定の状態のため開始できません。「設定」から入力してください。",
      "info",
      "need_config"
    );
    return;
  }

  const isFinishedPrepare = await IsFinishedPrepare();
  if (!isFinishedPrepare) {
    notification.showToastWithKey(
      "非プレイヤーデータの取得中...",
      "info",
      "prepare"
    );

    try {
      await Prepare();
    } catch (error) {
      notification.showToastWithKey(error, "error", "error");
      return;
    } finally {
      notification.removeToastWithKey("prepare");
    }
  }

  Ready();
}

window.onload = function () {
  main();
};
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
                disabled="{battle === undefined}"
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
          bind:battle="{battle}"
          bind:config="{config}"
          bind:averageFactors="{averageFactors}"
          bind:excludePlayerIDs="{excludePlayerIDs}"
        />
      </div>
    {/if}

    {#if currentPage === "config"}
      <ConfigPage
        on:onUpdateSuccess="{(event) => {
          notification.showToast(event.detail.message, 'success');
          UserConfig().then((result) => (config = result));
          notification.removeToastWithKey('need_config');
        }}"
        on:onUpdateFailure="{(event) =>
          notification.showToast(event.detail.message, 'error')}"
        on:onOpenDirectoryFailure="{(event) =>
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
