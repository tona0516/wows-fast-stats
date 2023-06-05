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

import "bootstrap-icons/font/bootstrap-icons.css";
import Navigation from "./Navigation.svelte";
import { Screenshot } from "./Screenshot.js";
import AlertPlayer from "./PageAlertPlayer.svelte";
import AddAlertPlayerModal from "./AddAlertPlayerModal.svelte";
import RemoveAlertPlayerModal from "./RemoveAlertPlayerModal.svelte";
import UpdateAlertPlayerModal from "./UpdateAlertPlayerModal.svelte";
import { get } from "svelte/store";
import {
  storedAlertPlayers,
  storedBattle,
  storedCurrentPage,
  storedExcludePlayerIDs,
  storedIsFirstScreenshot,
  storedSummaryResult,
  storedUserConfig,
} from "./stores.js";
import type { vo } from "wailsjs/go/models.js";
import { summary } from "./util.js";
import type { StatsPattern } from "./StatsPattern.js";

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

EventsOn("BATTLE_START", async () => {
  LogDebug("BATTLE_START");

  try {
    notification.removeToastWithKey("not_in_battle");
    notification.showToastWithKey("戦闘データの取得中...", "info", "battle");

    const start = new Date().getTime();

    const battle = await Battle();
    storedBattle.set(battle);
    recalculate(battle);

    const elapsed = (new Date().getTime() - start) / 1000;
    notification.showToast(`データ取得完了: ${elapsed}秒`, "success");
    notification.removeToastWithKey("error");
    notification.removeToastWithKey("battle");
  } catch (error) {
    notification.showToastWithKey(error, "error", "error");
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
      LogError(error.name + "," + error.message + "," + error.stack);
    }
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

async function recalculate(battle: vo.Battle) {
  const excludePlayerIDs = await ExcludePlayerIDs();
  storedExcludePlayerIDs.set(excludePlayerIDs);

  // TODO Refactoring (without "as")
  const summaryResult = summary(
    battle,
    excludePlayerIDs,
    userConfig.stats_pattern as StatsPattern
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

  try {
    const config = await UserConfig();
    storedUserConfig.set(config);

    if (!config.appid) {
      notification.showToastWithKey(
        "未設定の状態のため開始できません。「設定」から入力してください。",
        "info",
        "need_config"
      );
      return;
    }
  } catch (error) {
    notification.showToast(error, "error");
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
      on:ChangeStatsPattern="{() => recalculate(battle)}"
    />

    {#if currentPage === "main"}
      <div id="mainpage">
        <MainPage
          on:UpdateAlertPlayer="{(event) => showUpdateAlertPlayerModal(event)}"
          on:RemoveAlertPlayer="{(event) => showRemoveAlertPlayerModal(event)}"
          on:CheckPlayer="{() => recalculate(battle)}"
        />
      </div>
    {/if}

    {#if currentPage === "config"}
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
