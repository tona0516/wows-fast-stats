<script lang="ts">
  import AddAlertPlayerModal from "src/other_component/AddAlertPlayerModal.svelte";
  import MainPage from "src/page_component/MainPage.svelte";
  import ConfigPage from "src/page_component/ConfigPage.svelte";
  import AppInfoPage from "src/page_component/AppInfoPage.svelte";
  import AlertPlayerPage from "src/page_component/AlertPlayerPage.svelte";
  import UpdateAlertPlayerModal from "src/other_component/UpdateAlertPlayerModal.svelte";
  import RemoveAlertPlayerModal from "src/other_component/RemoveAlertPlayerModal.svelte";
  import Notification from "src/other_component/Notification.svelte";
  import Navigation from "src/other_component/Navigation.svelte";

  import "bootstrap-icons/font/bootstrap-icons.css";
  import {
    Battle,
    ExcludePlayerIDs,
    AlertPlayers,
    UserConfig,
    LatestRelease,
    StartWatching,
  } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";
  import {
    EventsOn,
    LogDebug,
    EventsEmit,
    BrowserOpenURL,
  } from "wailsjs/runtime/runtime";
  import { Screenshot } from "src/Screenshot";
  import { AppEvent, ToastKey, Page } from "./enums";
  import {
    storedSummaryResult,
    storedLogs,
    storedBattle,
    storedExcludePlayerIDs,
    storedAlertPlayers,
    storedUserConfig,
    storedCurrentPage,
  } from "./stores";
  import { summary } from "./util";

  let notification: Notification;
  let addAlertPlayerModal: AddAlertPlayerModal;
  let updateAlertPlayerModal: UpdateAlertPlayerModal;
  let removeAlertPlayerModal: RemoveAlertPlayerModal;

  let screenshot = new Screenshot();

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
        await screenshot.auto($storedBattle.meta);
      } catch (error) {
        notification.showToast(
          "スクリーンショットの自動保存に失敗しました。",
          "error"
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

  EventsOn(AppEvent.battleErr, (error: string) => {
    LogDebug(`EventsOn:${AppEvent.battleErr}`);

    notification.showToastWithKey(error, "error", ToastKey.error);
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

    try {
      StartWatching();
    } catch (error) {
      notification.showToast(error, "error");
    }
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
      {screenshot}
      on:Success={(event) =>
        notification.showToast(event.detail.message, "success")}
      on:Failure={(event) =>
        notification.showToast(event.detail.message, "error")}
    />

    {#if $storedCurrentPage === Page.Main}
      <MainPage
        on:UpdateAlertPlayer={(event) => showUpdateAlertPlayerModal(event)}
        on:RemoveAlertPlayer={(event) => showRemoveAlertPlayerModal(event)}
        on:CheckPlayer={async () =>
          storedExcludePlayerIDs.set(await ExcludePlayerIDs())}
      />
    {/if}

    {#if $storedCurrentPage === Page.Config}
      <ConfigPage
        on:UpdateSuccess={(event) => {
          notification.showToast(event.detail.message, "success");
          notification.removeToastWithKey(ToastKey.needConfig);
          if (!$storedBattle) {
            try {
              StartWatching();
            } catch (error) {
              notification.showToast(error, "error");
            }
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
