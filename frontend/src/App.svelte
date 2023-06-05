<script lang="ts">
import {
  UserConfig,
  Battle,
  ExcludePlayerIDs,
  Ready,
  LogError,
  AlertPlayers,
} from "../wailsjs/go/main/App.js";
import Notification from "./Notification.svelte";
import ConfigPage from "./PageConfig.svelte";
import MainPage from "./PageMain.svelte";
import { EventsOn, LogDebug } from "../wailsjs/runtime/runtime.js";
import AppInfo from "./PageAppInfo.svelte";
import { alertPlayers } from "./stores.js";

import "bootstrap-icons/font/bootstrap-icons.css";
import type { vo } from "wailsjs/go/models.js";
import { Summary, type SummaryResult } from "./Summary.js";
import Navigation from "./Navigation.svelte";
import type { Page } from "./Page.js";
import { Screenshot } from "./Screenshot.js";
import AlertPlayer from "./PageAlertPlayer.svelte";
import AddAlertPlayerModal from "./AddAlertPlayerModal.svelte";
import RemoveAlertPlayerModal from "./RemoveAlertPlayerModal.svelte";
import UpdateAlertPlayerModal from "./UpdateAlertPlayerModal.svelte";

let currentPage: Page;
let battle: vo.Battle;
let config: vo.UserConfig;
let summaryResult: SummaryResult;
let excludePlayerIDs: number[];
let notification: Notification;
let isFirstScreenshot: boolean;

let addAlertPlayerModal: AddAlertPlayerModal;
let updateAlertPlayerModal: UpdateAlertPlayerModal;
let removeAlertPlayerModal: RemoveAlertPlayerModal;

EventsOn("BATTLE_START", async () => {
  LogDebug("BATTLE_START");
  try {
    notification.removeToastWithKey("not_in_battle");
    notification.showToastWithKey("戦闘データの取得中...", "info", "battle");

    const start = new Date().getTime();

    battle = await Battle();
    excludePlayerIDs = await ExcludePlayerIDs();
    const summary = new Summary(battle);
    summaryResult = summary.calc(excludePlayerIDs);
    alertPlayers.set(await getAlertPlayers());

    const elapsed = (new Date().getTime() - start) / 1000;
    notification.showToast(`データ取得完了: ${elapsed}秒`, "success");
    notification.removeToastWithKey("error");

    if (config.save_screenshot) {
      const screenshot = new Screenshot(battle, isFirstScreenshot);
      screenshot
        .auto()
        .catch((error: Error) => {
          notification.showToast(
            "スクリーンショットの自動保存に失敗しました。",
            "error"
          );
          LogError(error.name + "," + error.message + "," + error.stack);
        })
        .finally(() => {
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
  LogDebug("BATTLE_END");
  notification.showToastWithKey(
    "戦闘中ではありません。開始時に自動的にリロードします。",
    "info",
    "not_in_battle"
  );
});

async function showAddAlertPlayerModal(_: CustomEvent<any>) {
  addAlertPlayerModal.toggle();
}

async function showUpdateAlertPlayerModal(event: CustomEvent<any>) {
  updateAlertPlayerModal.setTarget(event.detail.target);
  updateAlertPlayerModal.toggle();
}

async function showRemoveAlertPlayerModal(event: CustomEvent<any>) {
  removeAlertPlayerModal.setTarget(event.detail.target);
  removeAlertPlayerModal.toggle();
}

async function onSuccessAlertPlayerModal() {
  alertPlayers.set(await getAlertPlayers());
}

async function onFailureAlertPlayerModal(event: CustomEvent<any>) {
  notification.showToast(event.detail.message, "error");
}

async function getAlertPlayers(): Promise<vo.AlertPlayer[]> {
  try {
    return await AlertPlayers();
  } catch (error) {
    LogError(error);
    return [];
  }
}

async function main() {
  alertPlayers.set(await getAlertPlayers());

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
    <AddAlertPlayerModal
      bind:this="{addAlertPlayerModal}"
      on:Success="{onSuccessAlertPlayerModal}"
      on:Failure="{(event) => onFailureAlertPlayerModal(event)}"
    />

    <UpdateAlertPlayerModal
      bind:this="{updateAlertPlayerModal}"
      on:Success="{onSuccessAlertPlayerModal}"
      on:Failure="{(event) => onFailureAlertPlayerModal(event)}"
    />

    <RemoveAlertPlayerModal
      bind:this="{removeAlertPlayerModal}"
      on:Success="{onSuccessAlertPlayerModal}"
      on:Failure="{(event) => onFailureAlertPlayerModal(event)}"
    />

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
          bind:summaryResult="{summaryResult}"
          bind:excludePlayerIDs="{excludePlayerIDs}"
          on:UpdateAlertPlayer="{(event) => showUpdateAlertPlayerModal(event)}"
          on:RemoveAlertPlayer="{(event) => showRemoveAlertPlayerModal(event)}"
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

    {#if currentPage === "alert_player"}
      <AlertPlayer
        on:AddAlertPlayer="{(event) => showAddAlertPlayerModal(event)}"
        on:UpdateAlertPlayer="{(event) => showUpdateAlertPlayerModal(event)}"
        on:RemoveAlertPlayer="{(event) => showRemoveAlertPlayerModal(event)}"
        on:Failure="{(event) =>
          notification.showToast(event.detail.message, 'error')}"
      />
    {/if}

    <Notification bind:this="{notification}" />
  </div>
</main>
