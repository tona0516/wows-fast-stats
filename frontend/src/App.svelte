<script lang="ts">
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
    storedLogs,
  } from "./stores";
  import { summary } from "./util";
  import {
    AlertPlayers,
    Battle,
    ExcludePlayerIDs,
    LogErrorForFrontend,
    StartWatching,
    LatestRelease,
    UserConfig,
  } from "../wailsjs/go/main/App";
  import {
    BrowserOpenURL,
    EventsEmit,
    EventsOn,
    LogDebug,
  } from "../wailsjs/runtime/runtime";
  import type { domain } from "../wailsjs/go/models";
  import MainPage from "./page_component/MainPage.svelte";
  import ConfigPage from "./page_component/ConfigPage.svelte";
  import AppInfoPage from "./page_component/AppInfoPage.svelte";
  import AlertPlayerPage from "./page_component/AlertPlayerPage.svelte";
  import UpdateAlertPlayerModal from "./other_component/UpdateAlertPlayerModal.svelte";
  import RemoveAlertPlayerModal from "./other_component/RemoveAlertPlayerModal.svelte";
  import Notification from "./other_component/Notification.svelte";
  import Navigation from "./other_component/Navigation.svelte";
  import { AppEvent, Page, ToastKey } from "./enums";

  import "bootstrap-icons/font/bootstrap-icons.css";

  let notification: Notification;
  let addAlertPlayerModal: AddAlertPlayerModal;
  let updateAlertPlayerModal: UpdateAlertPlayerModal;
  let removeAlertPlayerModal: RemoveAlertPlayerModal;

  $: storedSummaryResult.set(
    summary($storedBattle, $storedExcludePlayerIDs, $storedUserConfig)
  );

  EventsOn(AppEvent.log, async (log: string) => {
    LogDebug(`EventsOn:${AppEvent.log}`);

    storedLogs.update((logs) => {
      logs.push(log);
      return logs;
    });
  });

  EventsOn(AppEvent.battleStart, async () => {
    LogDebug(`EventsOn:${AppEvent.battleStart}`);

    try {
      notification.removeToastWithKey(ToastKey.wait);
      notification.showToastWithKey(
        "戦闘データの取得中...",
        "info",
        ToastKey.fetching
      );

      const start = new Date().getTime();

      storedBattle.set(await Battle());
      storedExcludePlayerIDs.set(await ExcludePlayerIDs());

      const elapsed = (new Date().getTime() - start) / 1000;
      notification.showToast(`データ取得完了: ${elapsed}秒`, "success");
      notification.removeToastWithKey(ToastKey.error);
    } catch (error) {
      notification.showToastWithKey(error, "error", ToastKey.error);
    } finally {
      notification.removeToastWithKey(ToastKey.fetching);
    }

    if ($storedUserConfig.save_screenshot) {
      try {
        const screenshot = new Screenshot(
          $storedBattle,
          $storedIsFirstScreenshot
        );
        screenshot.auto();
        storedIsFirstScreenshot.set(false);
      } catch (error) {
        notification.showToast(
          "スクリーンショットの自動保存に失敗しました。",
          "error"
        );
        LogErrorForFrontend(
          error.name + "," + error.message + "," + error.stack
        );
      }
    }
  });

  EventsOn(AppEvent.battleEnd, () => {
    LogDebug(`EventsOn:${AppEvent.battleEnd}`);

    notification.showToastWithKey(
      "戦闘中ではありません。開始時に自動的にリロードします。",
      "info",
      ToastKey.wait
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

  async function main() {
    EventsEmit("ONLOAD");

    try {
      const players = await AlertPlayers();
      storedAlertPlayers.set(players);
    } catch (error) {
      notification.showToast(error, "error");
      return;
    }

    let config: domain.UserConfig;
    try {
      config = await UserConfig();
      storedUserConfig.set(config);
    } catch (error) {
      notification.showToast(error, "error");
      return;
    }

    if (config.notify_updatable) {
      try {
        const latestRelease = await LatestRelease();
        if (latestRelease.updatable) {
          notification.showToastWithKey(
            "新しいバージョンがあります: " +
              latestRelease.tag_name +
              "(クリックで開く)",
            "warning",
            ToastKey.updatable,
            () => BrowserOpenURL(latestRelease.html_url)
          );
        }
      } catch (error) {
        notification.showToast(error, "error");
      }
    }

    if (!config.appid) {
      notification.showToastWithKey(
        "未設定の状態のため開始できません。「設定」から入力してください。",
        "info",
        ToastKey.needConfig
      );
      return;
    }

    StartWatching();
  }

  window.onload = function () {
    main();
  };
</script>

<main>
  <div style="font-size: {$storedUserConfig.font_size};">
    <AddAlertPlayerModal
      bind:this={addAlertPlayerModal}
      on:Success={onSuccessAlertPlayerModal}
      on:Failure={(event) => onFailureAlertPlayerModal(event)}
    />

    <UpdateAlertPlayerModal
      bind:this={updateAlertPlayerModal}
      on:Success={onSuccessAlertPlayerModal}
      on:Failure={(event) => onFailureAlertPlayerModal(event)}
    />

    <RemoveAlertPlayerModal
      bind:this={removeAlertPlayerModal}
      on:Success={onSuccessAlertPlayerModal}
      on:Failure={(event) => onFailureAlertPlayerModal(event)}
    />

    <Navigation
      on:Success={(event) =>
        notification.showToast(event.detail.message, "success")}
      on:Failure={(event) =>
        notification.showToast(event.detail.message, "error")}
    />

    {#if $storedCurrentPage === Page.Main}
      <div id="mainpage">
        <MainPage
          on:UpdateAlertPlayer={(event) => showUpdateAlertPlayerModal(event)}
          on:RemoveAlertPlayer={(event) => showRemoveAlertPlayerModal(event)}
          on:CheckPlayer={async () =>
            storedExcludePlayerIDs.set(await ExcludePlayerIDs())}
        />
      </div>
    {/if}

    {#if $storedCurrentPage === Page.Config}
      <ConfigPage
        on:UpdateSuccess={(event) => {
          notification.showToast(event.detail.message, "success");
          notification.removeToastWithKey(ToastKey.needConfig);
          if (!$storedBattle) {
            StartWatching();
          }
        }}
        on:Failure={(event) =>
          notification.showToast(event.detail.message, "error")}
      />
    {/if}

    {#if $storedCurrentPage === Page.AppInfo}
      <AppInfoPage />
    {/if}

    {#if $storedCurrentPage === Page.AlertPlayer}
      <AlertPlayerPage
        on:AddAlertPlayer={(event) => showAddAlertPlayerModal(event)}
        on:UpdateAlertPlayer={(event) => showUpdateAlertPlayerModal(event)}
        on:RemoveAlertPlayer={(event) => showRemoveAlertPlayerModal(event)}
        on:Failure={(event) =>
          notification.showToast(event.detail.message, "error")}
      />
    {/if}

    <Notification bind:this={notification} />
  </div>
</main>
