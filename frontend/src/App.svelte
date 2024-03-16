<script lang="ts">
  import MainPage from "src/component/main/MainPage.svelte";
  import ConfigPage from "src/component/config/ConfigPage.svelte";

  import "bootstrap-icons/font/bootstrap-icons.css";
  import "charts.css";

  import {
    LatestRelease,
    LogError,
    MigrateIfNeeded,
    StartWatching,
  } from "wailsjs/go/main/App";
  import { EventsOn } from "wailsjs/runtime/runtime";
  import { storedConfig, storedRequiredConfigError } from "src/stores";
  import AlertModals from "src/component/modal/AlertModals.svelte";
  import { model } from "wailsjs/go/models";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkSpinner from "src/component/common/uikit/UkSpinner.svelte";
  import UkTab from "src/component/common/uikit/UkTab.svelte";
  import { FontSize } from "src/lib/FontSize";
  import { FetchProxy } from "./lib/FetchProxy";
  import { Notifier } from "./lib/Notifier";

  let modals: AlertModals;
  let mainPage: MainPage | undefined;
  let initialized = false;
  let updatableRelease: model.GHLatestRelease;

  $: {
    // @ts-ignore
    document.body.style.zoom = FontSize.getZoomRate($storedConfig);
  }

  EventsOn("BATTLE_START", () => mainPage?.fetchBattle());
  EventsOn("BATTLE_ERR", (error: string) => Notifier.failure(error));

  window.onunhandledrejection = (event) => {
    const message = "window.onunhandledrejection";
    const error = event.reason;
    if (error instanceof Error) {
      sendFronendError(message, error);
    } else {
      LogError(message, { error: JSON.stringify(error) });
    }
  };
  window.onerror = (_event, _source, _lineno, _colno, error) => {
    sendFronendError("window.onerror", error);
  };

  const sendFronendError = (message: string, error: Error | undefined) => {
    LogError(message, {
      "error.name": error?.name ?? "",
      "error.message": error?.message ?? "",
      "error.stack": error?.stack ?? "",
    });
  };

  const initialize = async (): Promise<model.UserConfigV2 | undefined> => {
    try {
      await MigrateIfNeeded();

      const config = await FetchProxy.getConfig();
      const requiredConfigError = await FetchProxy.validateRequiredConfig(
        config.install_path,
        config.appid,
      );
      await FetchProxy.getAlertPlayers();

      initialized = true;

      if (requiredConfigError.valid) {
        StartWatching();
      }

      return config;
    } catch (error) {
      Notifier.failure(error);
      return undefined;
    }
  };

  const notifyUpdate = async (config: model.UserConfigV2) => {
    if (!config.notify_updatable) return;

    try {
      const latestRelease = await LatestRelease();
      if (latestRelease.updatable) {
        updatableRelease = latestRelease;
      }
    } catch (error) {
      Notifier.failure(error);
      return;
    }
  };

  const main = async () => {
    const config = await initialize();
    if (!config) {
      return;
    }

    await notifyUpdate(config);
  };

  main();
</script>

<main>
  <AlertModals bind:this={modals} />

  {#if updatableRelease}
    <div class="uk-flex uk-flex-center uk-background-secondary">
      新しいバージョンがあります:
      <ExternalLink url={updatableRelease.html_url}>
        {updatableRelease.tag_name}
      </ExternalLink>
    </div>
  {/if}

  {#if initialized}
    {@const tabID = "page-tab"}
    <UkTab clazz="uk-margin-remove" id={tabID}>
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
    <ul id={tabID} class="uk-switcher">
      <li>
        <MainPage
          bind:this={mainPage}
          on:EditAlertPlayer={(e) => modals.showEdit(e.detail.target)}
          on:RemoveAlertPlayer={(e) => modals.showRemove(e.detail.target)}
        />
      </li>
      <li>
        <ConfigPage
          on:AddAlertPlayer={() => modals.showAdd()}
          on:EditAlertPlayer={(e) => modals.showEdit(e.detail.target)}
          on:RemoveAlertPlayer={(e) => modals.showRemove(e.detail.target)}
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
