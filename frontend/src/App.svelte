<script lang="ts">
  import "bootstrap-icons/font/bootstrap-icons.css";
  import "charts.css";
  import EditAlertPlayerModal from "@components/modals/EditAlertPlayerModal.svelte";
  import ExternalLink from "@components/ExternalLink.svelte";
  import SideMenu from "@components/SideMenu.svelte";
  import Toast from "@components/Toast.svelte";
  import {
    storedAlertPlayers,
    storedBattle,
    storedInstallPathError,
    storedZoomRate,
  } from "@libs/stores";
  import type { Page } from "@libs/types";
  import AlertPlayerPage from "@pages/AlertPlayerPage.svelte";
  import ConfigPage from "@pages/ConfigPage.svelte";
  import InfoPage from "@pages/InfoPage.svelte";
  import StatsPage from "@pages/StatsPage.svelte";
  import {
    InstallPath,
    ValidateInstallPath,
    AlertPlayers,
    SubscribeBattle,
    ShowMessageDialog,
    LogError,
  } from "@wails/go/main/App";
  import type { data } from "@wails/go/models";
  import { EventsOn } from "@wails/runtime/runtime";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import { TonakoManager } from "@libs/TonakoManager";
  import PlayerDetailModal from "@components/modals/PlayerDetailModal.svelte";
  import RemoveAlertPlayerModal from "@components/modals/RemoveAlertPlayerModal.svelte";
  import ShipDetailModal from "@components/modals/ShipDetailModal.svelte";

  let statsPage: StatsPage | undefined;
  let initialized = false;
  let updatableRelease: data.GHLatestRelease;

  let page: Page = "stats";

  $: {
    // @ts-ignore
    document.body.style.zoom = $storedZoomRate / 100;
  }

  onMount(() => {
    themeChange(false);
  });

  EventsOn("ALERT_PLAYERS_UPDATE", (players: data.AlertPlayer[]) =>
    storedAlertPlayers.set(players),
  );
  EventsOn("BATTLE_START", () => {
    TonakoManager.instance.setStartBattleState();
  });
  EventsOn("BATTLE_END", () => {
    TonakoManager.instance.setEndBattleState();
  });
  EventsOn("BATTLE_ERR", (message: string) => {
    TonakoManager.instance.setBattleErrorState(message);
  });
  EventsOn("BATTLE_FETCH_OTHERS", () => {
    TonakoManager.instance.setFetchOtherDataState();
  });
  EventsOn("BATTLE_FETCH_PLAYERS", () => {
    TonakoManager.instance.setFetchPlayerDataState();
  });
  EventsOn("BATTLE_FETCH_DONE", (battle: data.Battle) => {
    TonakoManager.instance.setHidden();
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

  const initialize = async (): Promise<void> => {
    localStorage.clear();

    try {
      const installPath = await InstallPath();
      const installPathError = await ValidateInstallPath(installPath);
      if (installPathError) {
        storedInstallPathError.set(installPathError);
      }

      const alertPlayers = await AlertPlayers();
      storedAlertPlayers.set(alertPlayers);

      initialized = true;

      if (!$storedInstallPathError) {
        SubscribeBattle();
      }
    } catch (error) {
      ShowMessageDialog(`初期化に失敗しました: ${error as string}`);
    }
  };

  // const notifyUpdate = async (config: data.UserConfigV2) => {
  //   return;

  //   if (!config.notify_updatable) return;

  //   try {
  //     const latestRelease = await LatestRelease();
  //     if (latestRelease.updatable) {
  //       updatableRelease = latestRelease;
  //     }
  //   } catch (error) {
  //     Notifier.failure(error);
  //     return;
  //   }
  // };

  const main = async () => {
    await initialize();
    // await notifyUpdate(config);
  };

  main();
</script>

<main>
  <div>
    <Toast />

    <EditAlertPlayerModal />
    <RemoveAlertPlayerModal />
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
