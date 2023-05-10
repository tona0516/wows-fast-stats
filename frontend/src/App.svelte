<script lang="ts">
import {
  UserConfig,
  Battle,
  ManualScreenshot,
  AutoScreenshot,
  ExcludePlayerIDs,
  Ready,
} from "../wailsjs/go/main/App.js";
import Notification from "./Notification.svelte";
import ConfigPage from "./PageConfig.svelte";
import MainPage from "./PageMain.svelte";
import { toPng } from "html-to-image";
import { EventsOn } from "../wailsjs/runtime/runtime.js";
import AppInfo from "./PageAppInfo.svelte";

import "bootstrap-icons/font/bootstrap-icons.css";
import type { vo } from "wailsjs/go/models.js";
import { Average, type AverageFactor } from "./Average.js";
import Navigation from "./Navigation.svelte";
import type { Page } from "./Page.js";

let currentPage: Page;
let battle: vo.Battle;
let config: vo.UserConfig;
let averageFactors: AverageFactor;
let excludePlayerIDs: number[];
let notification: Notification;

let firstScreenshot = true;
async function getScreenshotBase64(): Promise<[string, string]> {
  // Workaround: first screenshot cann't draw values in table.
  if (firstScreenshot) {
    await toPng(document.getElementById("mainpage"));
    firstScreenshot = false;
  }
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
    // TODO handling based on type
    const errStr = error as string;
    if (errStr.includes("Canceled")) {
      return;
    }
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

  Ready();
}

window.onload = function () {
  main();
};
</script>

<main>
  <div style="font-size: {config?.font_size || 'medium'};">
    <Navigation
      bind:currentPage="{currentPage}"
      bind:battle="{battle}"
      on:onScreenshot="{manualScreenshot}"
    />

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
          notification.removeToastWithKey('need_config');
          config = event.detail.config;
          if (!battle) {
            Ready();
          }
        }}"
        on:onUpdateFailure="{(event) =>
          notification.showToast(event.detail.message, 'error')}"
        on:onOpenDirectoryFailure="{(event) =>
          notification.showToast(event.detail.message, 'error')}"
      />
    {/if}

    {#if currentPage === "appinfo"}
      <AppInfo />
    {/if}

    <Notification bind:this="{notification}" />
  </div>
</main>
