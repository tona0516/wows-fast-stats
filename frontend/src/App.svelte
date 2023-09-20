<script lang="ts">
  import MainPage from "src/page_component/MainPage.svelte";
  import ConfigPage from "src/page_component/ConfigPage.svelte";

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
  import {
    EventsOn,
    EventsEmit,
    BrowserOpenURL,
  } from "wailsjs/runtime/runtime";
  import { Screenshot } from "src/lib/Screenshot";
  import { Notification } from "src/lib/Notification";
  import {
    storedSummary,
    storedLogs,
    storedBattle,
    storedExcludePlayerIDs,
    storedAlertPlayers,
    storedUserConfig,
    storedDefaultUserConfig,
  } from "src/stores";
  import {
    AppEvent,
    ScreenshotType,
    ZOOM_RATIO,
    type FontSize,
  } from "src/lib/types";
  import { Summary } from "src/lib/Summary";
  import Modals from "src/other_component/modal/Modals.svelte";
  import { domain } from "wailsjs/go/models";

  let modals: Modals;
  let mainPage: MainPage;

  let screenshot = new Screenshot();
  let notification = new Notification();
  let initialized = false;
  let isLoading = false;
  let updatableRelease: domain.GHLatestRelease;

  $: storedSummary.set(
    Summary.calculate(
      $storedBattle,
      $storedExcludePlayerIDs,
      $storedUserConfig,
    ),
  );
  $: isSatisfiedRequired =
    $storedUserConfig.install_path !== "" && $storedUserConfig.appid !== "";

  $: {
    const fontSize = $storedUserConfig.font_size as FontSize;
    const zoomRatio = ZOOM_RATIO.get(fontSize);
    if (zoomRatio) {
      document.body.style.zoom = zoomRatio;
    } else {
      document.body.style.zoom = 1.0;
    }
  }

  const PAGE_TAB_ID = "page-tab";

  EventsOn(AppEvent.LOG, async (log: string) => {
    storedLogs.update((logs) => {
      logs.push(log);
      return logs;
    });
  });

  EventsOn(AppEvent.BATTLE_START, async () => {
    try {
      isLoading = true;

      const start = new Date().getTime();

      const battle = await Battle();
      storedBattle.set(battle);
      const excludeIDs = await ExcludePlayerIDs();
      storedExcludePlayerIDs.set(excludeIDs);

      const summary = Summary.calculate(battle, excludeIDs, $storedUserConfig);
      storedSummary.set(summary);

      const elapsed = (new Date().getTime() - start) / 1000;

      notification.success(`データ取得完了: ${elapsed.toFixed(1)}秒`);

      if ($storedUserConfig.save_screenshot) {
        screenshot.take(ScreenshotType.auto, battle.meta);
      }
    } catch (error) {
      notification.failure(error);
    } finally {
      isLoading = false;
    }
  });

  EventsOn(AppEvent.BATTLE_ERR, (error: string) => {
    notification.failure(error);
  });

  async function main() {
    try {
      storedDefaultUserConfig.set(await DefaultUserConfig());
      storedUserConfig.set(await UserConfig());
      storedAlertPlayers.set(await AlertPlayers());
      initialized = true;

      if ($storedUserConfig.notify_updatable) {
        const latestRelease = await LatestRelease();
        if (latestRelease.updatable) {
          updatableRelease = latestRelease;
        }
      }

      if (isSatisfiedRequired) {
        StartWatching();
      }
    } catch (error) {
      notification.failure(error);
    }
  }

  window.onload = function () {
    EventsEmit("ONLOAD");
  };

  main();
</script>

<main>
  <Modals
    bind:this={modals}
    on:AlertPlayerUpdated={() => {
      AlertPlayers()
        .then((players) => storedAlertPlayers.set(players))
        .catch((error) => notification.failure(error));
    }}
    on:Failure={(event) => notification.failure(event.detail.message)}
  />

  {#if updatableRelease}
    <div class="uk-flex uk-flex-center uk-background-secondary">
      <span>新しいバージョンがあります: </span>
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a href="#" on:click={() => BrowserOpenURL(updatableRelease.html_url)}
        >{updatableRelease.tag_name}</a
      >
    </div>
  {/if}

  {#if initialized}
    <ul
      class="uk-margin-remove"
      uk-tab="connect: #{PAGE_TAB_ID}; animation: uk-animation-fade"
    >
      <li>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a href="#"><span uk-icon="icon: home"></span></a>
      </li>
      <li>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a href="#">
          <span uk-icon="icon: cog"></span>
          {#if !isSatisfiedRequired}
            <span class="uk-text-warning uk-text-small" uk-icon="icon: warning"
            ></span>
          {/if}
        </a>
      </li>
    </ul>
    <ul id={PAGE_TAB_ID} class="uk-switcher">
      <li>
        <MainPage
          bind:this={mainPage}
          {isLoading}
          {screenshot}
          on:EditAlertPlayer={(e) =>
            modals.showEditAlertPlayer(e.detail.target)}
          on:RemoveAlertPlayer={(e) =>
            modals.showRemoveAlertPlayer(e.detail.target)}
          on:CheckPlayer={() => {
            ExcludePlayerIDs().then((ids) => {
              storedExcludePlayerIDs.set(ids);
            });
          }}
          on:ScreenshotSaved={() =>
            notification.success("スクリーンショットを保存しました")}
          on:Failure={(event) => notification.failure(event.detail.message)}
        />
      </li>
      <li>
        <ConfigPage
          {isSatisfiedRequired}
          on:UpdateSuccess={async () => {
            notification.success("設定を更新しました");
            try {
              const config = await UserConfig();
              if (config.install_path && config.appid && !$storedBattle) {
                StartWatching();
              }
            } catch (error) {
              notification.failure(error);
            }
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
  {/if}
</main>
