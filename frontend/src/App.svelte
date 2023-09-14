<script lang="ts">
  import MainPage from "src/page_component/MainPage.svelte";
  import ConfigPage from "src/page_component/ConfigPage.svelte";
  import AppInfoPage from "src/page_component/AppInfoPage.svelte";
  import AlertPlayerPage from "src/page_component/AlertPlayerPage.svelte";
  import Notification from "src/other_component/Notification.svelte";
  import Navigation from "src/other_component/Navigation.svelte";

  import "bootstrap-icons/font/bootstrap-icons.css";
  import "charts.css";

  import {
    Battle,
    ExcludePlayerIDs,
    AlertPlayers,
    UserConfig,
    LatestRelease,
    StartWatching,
    DefaultUserConfig,
  } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";
  import {
    EventsOn,
    EventsEmit,
    BrowserOpenURL,
  } from "wailsjs/runtime/runtime";
  import { Screenshot } from "src/lib/Screenshot";
  import {
    storedSummary,
    storedLogs,
    storedBattle,
    storedExcludePlayerIDs,
    storedAlertPlayers,
    storedUserConfig,
    storedCurrentPage,
  } from "src/stores";
  import AddAlertPlayerModal from "src/other_component/AddAlertPlayerModal.svelte";
  import UpdateAlertPlayerModal from "src/other_component/UpdateAlertPlayerModal.svelte";
  import RemoveAlertPlayerModal from "src/other_component/RemoveAlertPlayerModal.svelte";
  import { AppEvent, Page, ScreenshotType, ToastKey } from "./lib/types";
  import { Summary } from "./lib/Summary";

  let defaultUserConfig: domain.UserConfig;
  let notification: Notification;
  let addAlertPlayerModal: AddAlertPlayerModal;
  let updateAlertPlayerModal: UpdateAlertPlayerModal;
  let removeAlertPlayerModal: RemoveAlertPlayerModal;

  let screenshot = new Screenshot();

  $: storedSummary.set(
    Summary.calculate(
      $storedBattle,
      $storedExcludePlayerIDs,
      $storedUserConfig,
    ),
  );

  EventsOn(AppEvent.LOG, async (log: string) => {
    storedLogs.update((logs) => {
      logs.push(log);
      return logs;
    });
  });

  EventsOn(AppEvent.BATTLE_START, async () => {
    try {
      notification.removeToastWithKey(ToastKey.WAIT);
      notification.showToastWithKey(
        "戦闘データの取得中...",
        "info",
        ToastKey.FETCHING,
      );

      const start = new Date().getTime();

      const battle = await Battle();
      storedBattle.set(battle);
      const excludeIDs = await ExcludePlayerIDs();
      storedExcludePlayerIDs.set(excludeIDs);

      const summary = Summary.calculate(battle, excludeIDs, $storedUserConfig);
      storedSummary.set(summary);

      const elapsed = (new Date().getTime() - start) / 1000;
      notification.showSuccessToast(`データ取得完了: ${elapsed}秒`);
      notification.removeToastWithKey(ToastKey.ERROR);

      if ($storedUserConfig.save_screenshot) {
        screenshot.take(ScreenshotType.auto, battle.meta);
      }
    } catch (error) {
      notification.showToastWithKey(error, "error", ToastKey.ERROR);
    } finally {
      notification.removeToastWithKey(ToastKey.FETCHING);
    }
  });

  EventsOn(AppEvent.BATTLE_END, () => {
    notification.showToastWithKey(
      "戦闘中ではありません。開始時に自動的にリロードします。",
      "info",
      ToastKey.WAIT,
    );
  });

  EventsOn(AppEvent.BATTLE_ERR, (error: string) => {
    notification.showToastWithKey(error, "error", ToastKey.ERROR);
  });

  const onAlertPlayerUpdated = async () => {
    try {
      const players = await AlertPlayers();
      storedAlertPlayers.set(players);
    } catch (error) {
      notification.showErrorToast(error);
    }
  };

  const onAlertPlayerUpdateFailed = (event: CustomEvent<any>) => {
    notification.showErrorToast(event.detail.message);
  };

  const showAddAlertPlayerModal = (_: CustomEvent<any>) => {
    addAlertPlayerModal.toggle();
  };

  const showUpdateAlertPlayerModal = (event: CustomEvent<any>) => {
    updateAlertPlayerModal.setTarget(event.detail.target);
    updateAlertPlayerModal.toggle();
  };

  const showRemoveAlertPlayerModal = (event: CustomEvent<any>) => {
    removeAlertPlayerModal.setTarget(event.detail.target);
    removeAlertPlayerModal.toggle();
  };

  async function main() {
    EventsEmit("ONLOAD");

    defaultUserConfig = await DefaultUserConfig();

    try {
      const players = await AlertPlayers();
      storedAlertPlayers.set(players);
    } catch (error) {
      notification.showErrorToast(error);
      return;
    }

    let config: domain.UserConfig;
    try {
      config = await UserConfig();
      storedUserConfig.set(config);
    } catch (error) {
      notification.showErrorToast(error);
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
            ToastKey.UPDATABLE,
            () => BrowserOpenURL(latestRelease.html_url),
          );
        }
      } catch (error) {
        notification.showErrorToast(error);
      }
    }

    if (!config.appid) {
      notification.showToastWithKey(
        "未設定の状態のため開始できません。「設定」から入力してください。",
        "info",
        ToastKey.NEED_CONFIG,
      );
      return;
    }

    try {
      StartWatching();
    } catch (error) {
      notification.showErrorToast(error);
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
      on:Success={onAlertPlayerUpdated}
      on:Failure={(event) => onAlertPlayerUpdateFailed(event)}
    />

    <UpdateAlertPlayerModal
      bind:this={updateAlertPlayerModal}
      on:Success={onAlertPlayerUpdated}
      on:Failure={(event) => onAlertPlayerUpdateFailed(event)}
    />

    <RemoveAlertPlayerModal
      bind:this={removeAlertPlayerModal}
      on:Success={onAlertPlayerUpdated}
      on:Failure={(event) => onAlertPlayerUpdateFailed(event)}
    />

    <Navigation
      {screenshot}
      on:Success={(event) =>
        notification.showSuccessToast(event.detail.message)}
      on:Failure={(event) => notification.showErrorToast(event.detail.message)}
    />

    {#if $storedCurrentPage === Page.MAIN}
      <MainPage
        on:UpdateAlertPlayer={(e) => showUpdateAlertPlayerModal(e)}
        on:RemoveAlertPlayer={(e) => showRemoveAlertPlayerModal(e)}
        on:CheckPlayer={async () =>
          storedExcludePlayerIDs.set(await ExcludePlayerIDs())}
      />
    {/if}

    {#if $storedCurrentPage === Page.CONFIG}
      <ConfigPage
        {defaultUserConfig}
        on:UpdateSuccess={(event) => {
          notification.showSuccessToast(event.detail.message);
          notification.removeToastWithKey(ToastKey.NEED_CONFIG);
          if (!$storedBattle) {
            try {
              StartWatching();
            } catch (error) {
              notification.showErrorToast(error);
            }
          }
        }}
        on:Failure={(event) =>
          notification.showErrorToast(event.detail.message)}
      />
    {/if}

    {#if $storedCurrentPage === Page.APPINFO}
      <AppInfoPage />
    {/if}

    {#if $storedCurrentPage === Page.ALERT_PLAYER}
      <AlertPlayerPage
        on:AddAlertPlayer={showAddAlertPlayerModal}
        on:UpdateAlertPlayer={(e) => showUpdateAlertPlayerModal(e)}
        on:RemoveAlertPlayer={(e) => showRemoveAlertPlayerModal(e)}
        on:Failure={(event) =>
          notification.showErrorToast(event.detail.message)}
      />
    {/if}

    <Notification bind:this={notification} />
  </div>
</main>
