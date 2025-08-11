<script lang="ts">
  import ConfigPage from "src/component/config/ConfigPage.svelte";
  import InfoPage from "src/component/info/InfoPage.svelte";
  import StatsPage from "src/component/stats/StatsPage.svelte";

  import "bootstrap-icons/font/bootstrap-icons.css";
  import "charts.css";

  import { FontSize } from "src/lib/FontSize";
  import {
    storedAlertPlayers,
    storedBattle,
    storedConfig,
    storedInstallPathError,
    storedTonako,
  } from "src/stores";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import {
    AlertPlayers,
    LatestRelease,
    LogError,
    MigrateIfNeeded,
    ShowMessageDialog,
    SubscribeBattle,
    UserConfig,
    ValidateInstallPath,
  } from "wailsjs/go/main/App";
  import type { data } from "wailsjs/go/models";
  import { EventsOn, LogInfo } from "wailsjs/runtime/runtime";
  import SideMenu from "./SideMenu.svelte";
  import ExternalLink from "./component/common/ExternalLink.svelte";
  import type { Page } from "./lib/types";
  import EditAlertPlayerModal from "./component/common/EditAlertPlayerModal.svelte";
  import DeleteAlertPlayerModal from "./component/common/DeleteAlertPlayerModal.svelte";
  import AlertPlayerPage from "./component/alert_player/AlertPlayerPage.svelte";
  import PlayerDetailModal from "./component/common/PlayerDetailModal.svelte";
  import Toast from "./component/common/Toast.svelte";
  import ShipDetailModal from "./component/common/ShipDetailModal.svelte";
  import { Tonako } from "./component/stats/internal/Tonako";

  let statsPage: StatsPage | undefined;
  let initialized = false;
  let updatableRelease: data.GHLatestRelease;

  let page: Page = "stats";

  $: {
    // @ts-ignore
    document.body.style.zoom = FontSize.getZoomRate($storedConfig);
  }

  onMount(() => {
    themeChange(false);
  });

  EventsOn("CONFIG_UPDATE", (config: data.UserConfigV2) =>
    storedConfig.set(config),
  );
  EventsOn("ALERT_PLAYERS_UPDATE", (players: data.AlertPlayer[]) =>
    storedAlertPlayers.set(players),
  );
  EventsOn("BATTLE_START", () => {
    storedBattle.set(undefined);
    storedTonako.set({
      message: "戦闘データを読み込み中",
      isLoading: true,
      tonako: Tonako.Standby,
    });
  });
  EventsOn("BATTLE_END", () => {
    storedTonako.set({
      message: "戦闘開始時に自動的にリロードします",
      isLoading: false,
      tonako: Tonako.Standby,
    });
  });
  EventsOn("BATTLE_ERR", (message: string) => {
    storedTonako.set({
      message: message,
      isLoading: false,
      tonako: Tonako.Sorry,
    });
  });
  EventsOn("BATTLE_FETCH_OTHERS", () => {
    storedTonako.set({
      message: "艦・マップ情報を取得中",
      isLoading: true,
      tonako: Tonako.Standby,
    });
  });
  EventsOn("BATTLE_FETCH_PLAYERS", () => {
    storedTonako.set({
      message: "プレイヤー情報を取得中",
      isLoading: true,
      tonako: Tonako.Standby,
    });
  });
  EventsOn("BATTLE_FETCH_DONE", (battle: data.Battle) => {
    storedTonako.set(undefined);
    storedBattle.set(battle);
  });

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
    // localStorage.clear();

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
        SubscribeBattle();
      }

      return config;
    } catch (error) {
      ShowMessageDialog(`初期化に失敗しました: ${error as string}`);
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
    <Toast />

    <EditAlertPlayerModal />
    <DeleteAlertPlayerModal />
    <PlayerDetailModal />
    <ShipDetailModal />

    <div class="flex divide-x-1 divide-neutral-500">
      <div class="flex-none z-10">
        <SideMenu bind:page />
      </div>

      <div class="flex-1 min-w-[1px] m-4">
        {#if updatableRelease}
          <div>
            新しいバージョンがあります:
            <ExternalLink url={updatableRelease.html_url}>
              {updatableRelease.tag_name}
            </ExternalLink>
          </div>
        {/if}

        {#if initialized}
          {#if page === "stats"}
            <StatsPage bind:this={statsPage} />
          {:else if page === "ap_config"}
            <AlertPlayerPage />
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
