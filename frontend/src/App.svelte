<script lang="ts">
  import "bootstrap-icons/font/bootstrap-icons.css";
  import "charts.css";
  import EditBlackListModal from "@components/modals/EditBlackListModal.svelte";
  import ExternalLink from "@components/ExternalLink.svelte";
  import SideMenu from "@components/SideMenu.svelte";
  import Toast from "@components/Toast.svelte";
  import {
    storedBlackList,
    storedBattle,
    storedInstallPathError,
    storedUserConfig,
  } from "@libs/stores";
  import type { Page } from "@libs/types";
  import BlackListPage from "@pages/BlackListPage.svelte";
  import ConfigPage from "@pages/ConfigPage.svelte";
  import InfoPage from "@pages/InfoPage.svelte";
  import StatsPage from "@pages/StatsPage.svelte";
  import {
    ValidateInstallPath,
    SubscribeBattle,
    ShowMessageDialog,
    LogError,
    GetBlackList,
    GetUserConfig,
  } from "@wails/go/main/App";
  import type { data, domain } from "@wails/go/models";
  import { EventsOn, LogInfo } from "@wails/runtime/runtime";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import { TonakoManager } from "@libs/TonakoManager";
  import PlayerDetailModal from "@components/modals/PlayerDetailModal.svelte";
  import RemoveBlackListModal from "@components/modals/RemoveBlackListModal.svelte";
  import ShipDetailModal from "@components/modals/ShipDetailModal.svelte";

  let statsPage: StatsPage | undefined;
  let initialized = false;
  let updatableRelease: data.GHLatestRelease;

  let page: Page = "stats";

  onMount(() => {
    themeChange(false);
  });

  EventsOn("BLACKLIST_UPDATE", (list: domain.BlackListItem[]) => {
    LogInfo("BLACKLIST_UPDATE");
    storedBlackList.set(list);
  });
  EventsOn("BATTLE_START", () => {
    LogInfo("BATTLE_START");
    TonakoManager.instance.setStartBattleState();
  });
  EventsOn("BATTLE_END", () => {
    LogInfo("BATTLE_END");
    TonakoManager.instance.setEndBattleState();
  });
  EventsOn("BATTLE_ERR", (message: string) => {
    LogInfo("BATTLE_ERR");
    TonakoManager.instance.setBattleErrorState(message);
  });
  EventsOn("BATTLE_FETCH_OTHERS", () => {
    LogInfo("BATTLE_FETCH_OTHERS");
    TonakoManager.instance.setFetchOtherDataState();
  });
  EventsOn("BATTLE_FETCH_PLAYERS", () => {
    LogInfo("BATTLE_FETCH_PLAYERS");
    TonakoManager.instance.setFetchPlayerDataState();
  });
  EventsOn("BATTLE_FETCH_DONE", (battle: data.Battle) => {
    LogInfo("BATTLE_FETCH_DONE");
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
    try {
      const userConfig = await GetUserConfig();
      storedUserConfig.set(userConfig);
      storedBlackList.set(await GetBlackList());

      const installPathError = await ValidateInstallPath(
        userConfig.install_path,
      );
      if (installPathError) {
        storedInstallPathError.set(installPathError);
      }

      initialized = true;

      if (!$storedInstallPathError) {
        LogInfo("call SubscribeBattle")
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

    <EditBlackListModal />
    <RemoveBlackListModal />
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
            <BlackListPage />
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
