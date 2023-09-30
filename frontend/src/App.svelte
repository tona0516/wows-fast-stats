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
  import { Notification } from "src/lib/Notification";
  import {
    storedExcludePlayerIDs,
    storedAlertPlayers,
    storedUserConfig,
    storedRequiredConfigError,
    storedLogs,
  } from "src/stores";
  import { AppEvent, ZOOM_RATIO, type FontSize } from "src/lib/types";
  import Modals from "src/component/modal/Modals.svelte";
  import { domain } from "wailsjs/go/models";
  import UkIcon from "./component/common/uikit/UkIcon.svelte";
  import ExternalLink from "./component/common/ExternalLink.svelte";
  import UkSpinner from "./component/common/uikit/UkSpinner.svelte";

  const PAGE_TAB_ID = "page-tab";

  let modals: Modals;
  let mainPage: MainPage | undefined;
  let notification = new Notification();
  let initialized = false;
  let updatableRelease: domain.GHLatestRelease;

  $: {
    const fontSize = $storedUserConfig.font_size as FontSize;
    const zoomRatio = ZOOM_RATIO.get(fontSize);
    if (zoomRatio) {
      // @ts-ignore
      document.body.style.zoom = zoomRatio;
    } else {
      // @ts-ignore
      document.body.style.zoom = 1.0;
    }
  }

  EventsOn(AppEvent.LOG, (log: string) =>
    storedLogs.update((logs) => {
      logs.push(log);
      return logs;
    }),
  );
  EventsOn(AppEvent.BATTLE_START, () => mainPage?.fetchBattle());
  EventsOn(AppEvent.BATTLE_ERR, (error: string) => notification.failure(error));

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
      notification.failure(error);
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
        notification.failure(error);
      }
    }}
    on:Failure={(event) => notification.failure(event.detail.message)}
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
    <ul
      class="uk-margin-remove"
      uk-tab="connect: #{PAGE_TAB_ID}; animation: uk-animation-fade"
    >
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
    </ul>
    <ul id={PAGE_TAB_ID} class="uk-switcher">
      <li>
        <MainPage
          bind:this={mainPage}
          on:FetchSuccess={(event) =>
            notification.success(event.detail.message)}
          on:EditAlertPlayer={(e) =>
            modals.showEditAlertPlayer(e.detail.target)}
          on:RemoveAlertPlayer={(e) =>
            modals.showRemoveAlertPlayer(e.detail.target)}
          on:CheckPlayer={async () => {
            storedExcludePlayerIDs.set(await ExcludePlayerIDs());
          }}
          on:ScreenshotSaved={() =>
            notification.success("スクリーンショットを保存しました")}
          on:Failure={(event) => notification.failure(event.detail.message)}
        />
      </li>
      <li>
        <ConfigPage
          on:UpdateRequired={() => {
            notification.success("設定を更新しました");
            StartWatching();
          }}
          on:UpdateSuccess={async () => {
            notification.success("設定を更新しました");
          }}
          on:AddAlertPlayer={() => modals.showAddAlertPlayer()}
          on:EditAlertPlayer={(e) =>
            modals.showEditAlertPlayer(e.detail.target)}
          on:RemoveAlertPlayer={(e) =>
            modals.showRemoveAlertPlayer(e.detail.target)}
          on:Failure={(event) => notification.failure(event.detail.message)}
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
