<script lang="ts">
  import "bootstrap-icons/font/bootstrap-icons.css";
  import "charts.css";
  import SideMenu from "@components/SideMenu.svelte";
  import Toast from "@components/Toast.svelte";
  import { storedBattle, storedPref } from "@libs/stores";
  import type { Page } from "@libs/types";
  import PrefPage from "@pages/PrefPage.svelte";
  import InfoPage from "@pages/InfoPage.svelte";
  import StatsPage from "@pages/StatsPage.svelte";
  import {
    ShowMessageDialog,
    LoadPref,
    StartPollingMatch,
    Prefetch,
  } from "@wails/go/main/App";
  import type { core } from "@wails/go/models";
  import { EventsOn, LogInfo } from "@wails/runtime/runtime";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import { TonakoManager } from "@libs/TonakoManager";
  import ShipDetailModal from "@components/modals/ShipDetailModal.svelte";
  import PlayerDetailModal from "@components/modals/PlayerDetailModal.svelte";

  let statsPage: StatsPage | undefined;
  let initialized = false;

  let page: Page = "stats";

  $: {
    // @ts-ignore
    document.body.style.zoom = ($storedPref?.zoom_rate || 1.0) / 100;
  }

  onMount(() => {
    themeChange(false);
  });

  EventsOn("NEED_INITIAL_SETTING", () => {
    LogInfo("NEED_INITIAL_SETTING");
    TonakoManager.getInstance.setNeedInitialSettingState();
  });
  EventsOn("POLLING_START", () => {
    LogInfo("POLLING_START");
    TonakoManager.getInstance.setPollingStartState();
  });
  EventsOn("BATTLE_START", () => {
    LogInfo("BATTLE_START");
    TonakoManager.getInstance.setStartBattleState();
  });
  EventsOn("BATTLE_ERR", (message: string) => {
    LogInfo("BATTLE_ERR");
    TonakoManager.getInstance.setBattleErrorState(message);
  });
  EventsOn("BATTLE_FETCH_OTHERS", () => {
    LogInfo("BATTLE_FETCH_OTHERS");
    TonakoManager.getInstance.setFetchOtherDataState();
  });
  EventsOn("BATTLE_FETCH_PLAYERS", () => {
    LogInfo("BATTLE_FETCH_PLAYERS");
    TonakoManager.getInstance.setFetchPlayerDataState();
  });
  EventsOn("BATTLE_FETCH_DONE", (battle: core.Battle) => {
    LogInfo("BATTLE_FETCH_DONE");
    TonakoManager.getInstance.setHidden();
    storedBattle.set(battle);
  });

  const initialize = async (): Promise<void> => {
    try {
      const pref = await LoadPref();
      storedPref.set(pref);

      await Prefetch();

      StartPollingMatch();

      initialized = true;
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

    <div class="flex divide-x divide-neutral-500">
      <div class="flex-none z-10">
        <SideMenu bind:page />
      </div>

      <div class="flex-1 min-w-px m-4">
        {#if initialized}
          {#if page === "stats"}
            <StatsPage bind:this={statsPage} />
          {:else if page === "pref"}
            <PrefPage />
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
