<script lang="ts">
  import "bootstrap-icons/font/bootstrap-icons.css";
  import "charts.css";
  import SideMenu from "@components/SideMenu.svelte";
  import Toast from "@components/Toast.svelte";
  import {
    storedBattle,
    storedInstallPathError,
    storedUserConfig,
  } from "@libs/stores";
  import type { Page } from "@libs/types";
  import ConfigPage from "@pages/ConfigPage.svelte";
  import InfoPage from "@pages/InfoPage.svelte";
  import StatsPage from "@pages/StatsPage.svelte";
  import {
    ValidateInstallPath,
    SubscribeBattle,
    ShowMessageDialog,
    GetUserConfig,
  } from "@wails/go/main/App";
  import type { data } from "@wails/go/models";
  import { EventsOn, LogInfo } from "@wails/runtime/runtime";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import { TonakoManager } from "@libs/TonakoManager";
  import ShipDetailModal from "@components/modals/ShipDetailModal.svelte";
  import PlayerDetailModal from "@components/modals/PlayerDetailModal.svelte";

  let statsPage: StatsPage | undefined;
  let initialized = false;

  let page: Page = "stats";

  onMount(() => {
    themeChange(false);
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

  const initialize = async (): Promise<void> => {
    try {
      const userConfig = await GetUserConfig();
      storedUserConfig.set(userConfig);

      const installPathError = await ValidateInstallPath(
        userConfig.install_path,
      );
      if (installPathError) {
        storedInstallPathError.set(installPathError);
      }

      initialized = true;

      if (!$storedInstallPathError) {
        LogInfo("call SubscribeBattle");
        SubscribeBattle();
      }
    } catch (error) {
      ShowMessageDialog(`初期化に失敗しました: ${error as string}`);
    }
  };

  const main = async () => {
    await initialize();
  };

  main();
</script>

<main>
  <div>
    <Toast />
    <PlayerDetailModal />
    <ShipDetailModal />

    <div class="flex divide-x-1 divide-neutral-500">
      <div class="flex-none z-10">
        <SideMenu bind:page />
      </div>

      <div class="flex-1 min-w-[1px] m-4">
        {#if initialized}
          {#if page === "stats"}
            <StatsPage bind:this={statsPage} />
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
