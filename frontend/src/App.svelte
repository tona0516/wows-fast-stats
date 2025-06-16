<script lang="ts">
import ConfigPage from "src/component/config/ConfigPage.svelte";
import InfoPage from "src/component/info/InfoPage.svelte";
import MainPage from "src/component/main/MainPage.svelte";

import "bootstrap-icons/font/bootstrap-icons.css";
import "charts.css";

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
import ExternalLink from "./component/common/ExternalLink.svelte";
import { Notifier } from "./lib/Notifier";

let mainPage: MainPage | undefined;
let initialized = false;
let updatableRelease: data.GHLatestRelease;

type Page = "main" | "config" | "info";
let page: Page = "main";

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
  return;

  // if (!config.notify_updatable) return;

  // try {
  //   const latestRelease = await LatestRelease();
  //   if (latestRelease.updatable) {
  //     updatableRelease = latestRelease;
  //   }
  // } catch (error) {
  //   Notifier.failure(error);
  //   return;
  // }
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
  <div>
    <div class="flex">
      <div class="flex-none">
        <ul class="menu rounded-box sticky top-0">
          <li>
            <a on:click={() => (page = "main")}>
              <i class="bi bi-house"></i>
            </a>
          </li>
          <li>
            <a on:click={() => (page = "config")}>
              <i class="bi bi-gear"></i>
            </a>
          </li>
          <li>
            <a on:click={() => (page = "info")}
              ><i class="bi bi-info-circle"></i></a
            >
          </li>
        </ul>
      </div>

      <div class="flex-1 min-w-[1px]">
        {#if updatableRelease}
          <div>
            新しいバージョンがあります:
            <ExternalLink url={updatableRelease.html_url}>
              {updatableRelease.tag_name}
            </ExternalLink>
          </div>
        {/if}

        {#if initialized}
          {#if page === "main"}
            <MainPage bind:this={mainPage} />
          {:else if page === "config"}
            <ConfigPage />
          {:else if page === "info"}
            <InfoPage />
          {/if}
        {:else}
          <div
            class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-20"
          >
            <span class="loading loading-spinner"></span>
          </div>
        {/if}
      </div>
    </div>
  </div>
</main>
