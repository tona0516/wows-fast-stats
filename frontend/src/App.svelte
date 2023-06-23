<script lang="ts">
import { get } from "svelte/store";
import AddAlertPlayerModal from "./other_component/AddAlertPlayerModal.svelte";
import { Screenshot } from "./Screenshot";
import {
  storedBattle,
  storedUserConfig,
  storedIsFirstScreenshot,
  storedCurrentPage,
  storedAlertPlayers,
  storedExcludePlayerIDs,
  storedSummaryResult,
} from "./stores";
import { summary } from "./util";
import {
  AlertPlayers,
  Battle,
  ExcludePlayerIDs,
  LogErrorForFrontend,
  Ready,
  UserConfig,
} from "../wailsjs/go/main/App";
import { EventsOn, LogDebug } from "../wailsjs/runtime/runtime";
import type { vo } from "../wailsjs/go/models";
import MainPage from "./page_component/MainPage.svelte";
import ConfigPage from "./page_component/ConfigPage.svelte";
import AppInfoPage from "./page_component/AppInfoPage.svelte";
import AlertPlayerPage from "./page_component/AlertPlayerPage.svelte";
import UpdateAlertPlayerModal from "./other_component/UpdateAlertPlayerModal.svelte";
import RemoveAlertPlayerModal from "./other_component/RemoveAlertPlayerModal.svelte";
import Notification from "./other_component/Notification.svelte";
import Navigation from "./other_component/Navigation.svelte";
import { Page } from "./enums";

import "bootstrap-icons/font/bootstrap-icons.css";

const TOAST_NEED_CONFIG = "need_config";
const TOAST_WAIT = "wait";
const TOAST_FETCHING = "fetching";
const TOAST_ERROR = "error";

// Note: see watcher.go
const EVENT_BATTLE_START = "BATTLE_START"
const EVENT_BATTLE_END = "BATTLE_END"

let notification: Notification;
let addAlertPlayerModal: AddAlertPlayerModal;
let updateAlertPlayerModal: UpdateAlertPlayerModal;
let removeAlertPlayerModal: RemoveAlertPlayerModal;

let battle = get(storedBattle);
storedBattle.subscribe((it) => (battle = it));

let userConfig = get(storedUserConfig);
storedUserConfig.subscribe((it) => (userConfig = it));

let isFirstScreenshot = get(storedIsFirstScreenshot);
storedIsFirstScreenshot.subscribe((it) => (isFirstScreenshot = it));

let currentPage = get(storedCurrentPage);
storedCurrentPage.subscribe((it) => (currentPage = it));

EventsOn(EVENT_BATTLE_START, async () => {
  LogDebug(EVENT_BATTLE_START);

  try {
    notification.removeToastWithKey(TOAST_WAIT);
    notification.showToastWithKey("戦闘データの取得中...", "info", TOAST_FETCHING);

    const start = new Date().getTime();

    const battle = await Battle();
    storedBattle.set(battle);
    updateSummary(battle);

    const elapsed = (new Date().getTime() - start) / 1000;
    notification.showToast(`データ取得完了: ${elapsed}秒`, "success");
    notification.removeToastWithKey(TOAST_ERROR);
  } catch (error) {
    notification.showToastWithKey(error, "error", TOAST_ERROR);
  } finally {
    notification.removeToastWithKey(TOAST_FETCHING);
  }

  if (userConfig.save_screenshot) {
    try {
      const screenshot = new Screenshot(battle, isFirstScreenshot);
      screenshot.auto();
      storedIsFirstScreenshot.set(false);
    } catch (error) {
      notification.showToast(
        "スクリーンショットの自動保存に失敗しました。",
        "error"
      );
      LogErrorForFrontend(error.name + "," + error.message + "," + error.stack);
    }
  }
});

EventsOn(EVENT_BATTLE_END, () => {
  LogDebug(EVENT_BATTLE_END);

  notification.showToastWithKey(
    "戦闘中ではありません。開始時に自動的にリロードします。",
    "info",
    TOAST_WAIT
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
  try {
    const players = await AlertPlayers();
    storedAlertPlayers.set(players);
  } catch (error) {
    notification.showToast(error, "error");
  }
}

async function onFailureAlertPlayerModal(event: CustomEvent<any>) {
  notification.showToast(event.detail.message, "error");
}

async function updateSummary(battle: vo.Battle) {
  const excludePlayerIDs = await ExcludePlayerIDs();
  storedExcludePlayerIDs.set(excludePlayerIDs);

  const summaryResult = summary(
    battle,
    excludePlayerIDs,
    userConfig,
  );
  storedSummaryResult.set(summaryResult);
}

async function main() {
  try {
    const players = await AlertPlayers();
    storedAlertPlayers.set(players);
  } catch (error) {
    notification.showToast(error, "error");
    return;
  }

  let config: vo.UserConfig
  try {
    config = await UserConfig();
    storedUserConfig.set(config);
  } catch (error) {
    notification.showToast(error, "error");
    return;
  }

  if (!config.appid) {
      notification.showToastWithKey(
        "未設定の状態のため開始できません。「設定」から入力してください。",
        "info",
        TOAST_NEED_CONFIG
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
  <div style="font-size: {userConfig.font_size};">
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
      on:Success="{(event) =>
        notification.showToast(event.detail.message, 'success')}"
      on:Failure="{(event) =>
        notification.showToast(event.detail.message, 'error')}"
      on:ChangeStatsPattern="{() => updateSummary(battle)}"
    />

    {#if currentPage === Page.Main}
      <div id="mainpage">
        <MainPage
          on:UpdateAlertPlayer="{(event) => showUpdateAlertPlayerModal(event)}"
          on:RemoveAlertPlayer="{(event) => showRemoveAlertPlayerModal(event)}"
          on:CheckPlayer="{() => updateSummary(battle)}"
        />
      </div>
    {/if}

    {#if currentPage === Page.Config}
      <ConfigPage
        on:UpdateSuccess="{(event) => {
          notification.showToast(event.detail.message, 'success');
          notification.removeToastWithKey('need_config');
          if (!battle) {
            Ready();
          }
        }}"
        on:Failure="{(event) =>
          notification.showToast(event.detail.message, 'error')}"
      />
    {/if}

    {#if currentPage === Page.AppInfo}
      <AppInfoPage />
    {/if}

    {#if currentPage === Page.AlertPlayer}
      <AlertPlayerPage
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
