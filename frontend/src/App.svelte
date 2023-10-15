<script lang="ts">
  import MainPage from "src/component/main/MainPage.svelte";
  import ConfigPage from "src/component/config/ConfigPage.svelte";

  import "bootstrap-icons/font/bootstrap-icons.css";
  import "charts.css";

  import {
    ExcludePlayerIDs,
    AlertPlayers,
    UserConfig,
    LatestRelease,
    StartWatching,
    ValidateRequiredConfig,
  } from "wailsjs/go/main/App";
  import { EventsOn, EventsEmit } from "wailsjs/runtime/runtime";
  import {
    storedExcludePlayerIDs,
    storedAlertPlayers,
    storedUserConfig,
    storedRequiredConfigError,
    storedLogs,
  } from "src/stores";
  import Modals from "src/component/modal/AlertModals.svelte";
  import { domain } from "wailsjs/go/models";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkSpinner from "src/component/common/uikit/UkSpinner.svelte";
  import UkTab from "src/component/common/uikit/UkTab.svelte";
  import { FontSize } from "src/lib/FontSize";
  import { Notification } from "src/lib/Notification";

  // Note: see watcher.go
  enum AppEvent {
    BATTLE_START = "BATTLE_START",
    BATTLE_ERR = "BATTLE_ERR",
    LOG = "LOG",
    ONLOAD = "ONLOAD",
  }

  const PAGE_TAB_ID = "page-tab";

  let modals: Modals;
  let mainPage: MainPage | undefined;
  let initialized = false;
  let updatableRelease: domain.GHLatestRelease;

  $: {
    // @ts-ignore
    document.body.style.zoom = FontSize.getZoomRate($storedUserConfig);
  }

  EventsOn(AppEvent.LOG, (log: string) =>
    storedLogs.update((logs) => {
      logs.push(log);
      return logs;
    }),
  );
  EventsOn(AppEvent.BATTLE_START, () => mainPage?.fetchBattle());
  EventsOn(AppEvent.BATTLE_ERR, (error: string) => Notification.failure(error));

  async function main() {
    try {
      const userConfig = await UserConfig();
      const requiredConfigError = await ValidateRequiredConfig(
        userConfig.install_path,
        userConfig.appid,
      );

      storedUserConfig.set(userConfig);
      storedRequiredConfigError.set(requiredConfigError);
      storedAlertPlayers.set(await AlertPlayers());

      initialized = true;

      if (userConfig.notify_updatable) {
        const latestRelease = await LatestRelease();
        if (latestRelease.updatable) {
          updatableRelease = latestRelease;
        }
      }

      if (requiredConfigError.valid) {
        StartWatching();
      }
    } catch (error) {
      Notification.failure(error);
    }
  }

  window.onload = function () {
    EventsEmit(AppEvent.ONLOAD);
  };

  main();
</script>

<main>
  <Modals
    bind:this={modals}
    on:AlertPlayerUpdated={async () => {
      try {
        storedAlertPlayers.set(await AlertPlayers());
      } catch (error) {
        Notification.failure(error);
      }
    }}
    on:Failure={(event) => Notification.failure(event.detail.message)}
  />

  {#if updatableRelease}
    <div class="uk-flex uk-flex-center uk-background-secondary">
      新しいバージョンがあります:
      <ExternalLink url={updatableRelease.html_url}>
        {updatableRelease.tag_name}
      </ExternalLink>
    </div>
  {/if}

  {#if initialized}
    <UkTab clazz="uk-margin-remove" id={PAGE_TAB_ID}>
      <li>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a href="#"><UkIcon name="home" /></a>
      </li>
      <li>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a href="#">
          <UkIcon name="cog" />
          {#if !$storedRequiredConfigError.valid}
            <span class="uk-text-warning uk-text-small">
              <UkIcon name="warning" />
            </span>
          {/if}
        </a>
      </li>
    </UkTab>
    <ul id={PAGE_TAB_ID} class="uk-switcher">
      <li>
        <MainPage
          bind:this={mainPage}
          on:FetchSuccess={(event) =>
            Notification.success(event.detail.message)}
          on:EditAlertPlayer={(e) =>
            modals.showEditAlertPlayer(e.detail.target)}
          on:RemoveAlertPlayer={(e) =>
            modals.showRemoveAlertPlayer(e.detail.target)}
          on:CheckPlayer={async () => {
            storedExcludePlayerIDs.set(await ExcludePlayerIDs());
          }}
          on:ScreenshotSaved={() =>
            Notification.success("スクリーンショットを保存しました")}
          on:Failure={(event) => Notification.failure(event.detail.message)}
        />
      </li>
      <li>
        <ConfigPage
          on:UpdateRequired={() => {
            Notification.success("設定を更新しました");
            StartWatching();
          }}
          on:UpdateSuccess={async () => {
            Notification.success("設定を更新しました");
          }}
          on:AddAlertPlayer={() => modals.showAddAlertPlayer()}
          on:EditAlertPlayer={(e) =>
            modals.showEditAlertPlayer(e.detail.target)}
          on:RemoveAlertPlayer={(e) =>
            modals.showRemoveAlertPlayer(e.detail.target)}
          on:Failure={(event) => Notification.failure(event.detail.message)}
        />
      </li>
    </ul>
  {:else}
    <div class="uk-overlay-default">
      <div class="uk-position-center">
        <UkSpinner />
      </div>
    </div>
  {/if}
</main>
