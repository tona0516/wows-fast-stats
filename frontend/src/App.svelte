<script lang="ts">
import {
  UserConfig,
  Battle,
  ExcludePlayerIDs,
  Ready,
} from "../wailsjs/go/main/App.js";
import Notification from "./Notification.svelte";
import ConfigPage from "./PageConfig.svelte";
import MainPage from "./PageMain.svelte";
import { EventsOn } from "../wailsjs/runtime/runtime.js";
import AppInfo from "./PageAppInfo.svelte";

import "bootstrap-icons/font/bootstrap-icons.css";
import type { vo } from "wailsjs/go/models.js";
import { Average, type AverageFactor } from "./Average.js";
import Navigation from "./Navigation.svelte";
import type { Page } from "./Page.js";
import { Screenshot } from "./Screenshot.js";

let currentPage: Page;
let battle: vo.Battle;
let config: vo.UserConfig;
let averageFactors: AverageFactor;
let excludePlayerIDs: number[];
let notification: Notification;
let isFirstScreenshot: boolean;

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
      const screenshot = new Screenshot(battle, isFirstScreenshot);
      screenshot.auto().finally(() => {
        isFirstScreenshot = false;
      });
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
      bind:config="{config}"
      bind:currentPage="{currentPage}"
      bind:battle="{battle}"
      bind:isFirstScreenshot="{isFirstScreenshot}"
      on:onScreenshotSuccess="{(event) =>
        notification.showToast(event.detail.message, 'success')}"
      on:onScreenshotFailure="{(event) =>
        notification.showToast(event.detail.message, 'error')}"
    />

    {#if currentPage === "main"}
      <div id="mainpage">
        <MainPage
          bind:config="{config}"
          bind:battle="{battle}"
          bind:averageFactors="{averageFactors}"
          bind:excludePlayerIDs="{excludePlayerIDs}"
        />
      </div>
    {/if}

    {#if currentPage === "config"}
      <ConfigPage
        bind:config="{config}"
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
