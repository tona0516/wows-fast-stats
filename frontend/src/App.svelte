<script lang="ts">
import ConfigPage from "src/component/config/ConfigPage.svelte";
import InfoPage from "src/component/info/InfoPage.svelte";
import MainPage from "src/component/main/MainPage.svelte";

import "bootstrap-icons/font/bootstrap-icons.css";
import "charts.css";

import ExternalLink from "src/component/common/ExternalLink.svelte";
import UkSpinner from "src/component/common/uikit/UkSpinner.svelte";
import AlertModals from "src/component/modal/AlertModals.svelte";
import { FontSize } from "src/lib/FontSize";
import {
  storedAlertPlayers,
  storedConfig,
  storedInstallPathError,
  storedLogs,
} from "src/stores";
import {
  AlertPlayers,
  LatestRelease,
  LogError,
  MigrateIfNeeded,
  StartWatching,
  UserConfig,
  ValidateInstallPath,
} from "wailsjs/go/main/App";
import { data } from "wailsjs/go/models";
import { EventsOn } from "wailsjs/runtime/runtime";
import { Notifier } from "./lib/Notifier";

let modals: AlertModals;
let mainPage: MainPage | undefined;
let initialized = false;
let updatableRelease: data.GHLatestRelease;

$: {
  // @ts-ignore
  document.body.style.zoom = FontSize.getZoomRate($storedConfig);
}

EventsOn("BATTLE_START", () => mainPage?.fetchBattle());
EventsOn("BATTLE_ERR", (error: string) => Notifier.failure(error));
EventsOn("CONFIG_UPDATE", (config: data.UserConfigV2) =>
  storedConfig.set(config),
);
EventsOn("ALERT_PLAYERS_UPDATE", (players: data.AlertPlayer[]) =>
  storedAlertPlayers.set(players),
);
EventsOn("LOG", (log: string) =>
  storedLogs.update((logs) => {
    logs.push(log);
    return logs;
  }),
);

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

const initialize = async (): Promise<data.UserConfigV2 | undefined> => {
  try {
    await MigrateIfNeeded();

    const config = await UserConfig();
    storedConfig.set(config);

    const installPathError = await ValidateInstallPath(config.install_path);
    if (installPathError) {
      storedInstallPathError.set(installPathError);
    }

    const alertPlayers = await AlertPlayers();
    storedAlertPlayers.set(alertPlayers);

    initialized = true;

    if (!$storedInstallPathError) {
      StartWatching();
    }

    return config;
  } catch (error) {
    Notifier.failure(error);
    return undefined;
  }
};

const notifyUpdate = async (config: data.UserConfigV2) => {
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
    <div>
      新しいバージョンがあります:
      <ExternalLink url={updatableRelease.html_url}>
        {updatableRelease.tag_name}
      </ExternalLink>
    </div>
  {/if}

  {#if initialized}
    {@const tabID = "page-tab"}
    <div class="tabs tabs-border">
      <label class="tab">
        <input type="radio" name={tabID} checked={true} />
        ホーム
      </label>
      <div class="tab-content">
        <MainPage
          bind:this={mainPage}
          on:EditAlertPlayer={(e) => modals.showEdit(e.detail.target)}
          on:RemoveAlertPlayer={(e) => modals.showRemove(e.detail.target)}
        />
      </div>

      <label class="tab">
        <input type="radio" name={tabID} />
        設定
      </label>
      <div class="tab-content">
        <ConfigPage
          on:AddAlertPlayer={() => modals.showAdd()}
          on:EditAlertPlayer={(e) => modals.showEdit(e.detail.target)}
          on:RemoveAlertPlayer={(e) => modals.showRemove(e.detail.target)}
        />
      </div>

      <label class="tab">
        <input type="radio" name={tabID} />
        アプリについて
      </label>
      <div class="tab-content">
        <InfoPage />
      </div>
    </div>
  {:else}
    <div class="uk-overlay-default">
      <div class="uk-position-center">
        <UkSpinner />
      </div>
    </div>
  {/if}
</main>
