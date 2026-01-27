<script lang="ts">
  import "bootstrap-icons/font/bootstrap-icons.css";
  import "charts.css";
  import SideMenu from "@components/SideMenu.svelte";
  import Toast from "@components/Toast.svelte";
  import {
    storedDisplayPref,
    storedBattle,
    storedGameClientPathError,
    storedGameClientPath,
  } from "@libs/stores";
  import type { Page } from "@libs/types";
  import PrefPage from "@pages/PrefPage.svelte";
  import InfoPage from "@pages/InfoPage.svelte";
  import StatsPage from "@pages/StatsPage.svelte";
  import {
    StartPollingMatch,
    Prefetch,
    FetchBattle,
    LoadDisplayPref,
    LoadGameClientPath,
  } from "@wails/go/main/App";
  import type { core } from "@wails/go/models";
  import { EventsOn } from "@wails/runtime/runtime";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import { TonakoManager } from "@libs/TonakoManager";
  import ShipDetailModal from "@components/modals/ShipDetailModal.svelte";
  import PlayerDetailModal from "@components/modals/PlayerDetailModal.svelte";
    import { DEFAULT_DISPLAY_PREF, type DisplayPref } from "@libs/DisplayPref";

  let page: Page = "stats";

  $: {
    // @ts-ignore
    document.body.style.zoom = ($storedDisplayPref?.zoomRate || 1.0) / 100;
  }

  onMount(() => {
    themeChange(false);
  });

  EventsOn("ON_GAME_CLIENT_PATH_REQUIRED", () => {
    const message = "ゲームクライアントパスを設定してください";
    storedGameClientPathError.set(message);
    TonakoManager.getInstance.setPromoteState(message);
  });
  EventsOn("ON_START_POLLING", () => {
    TonakoManager.getInstance.setStandbyState("戦闘開始時に自動的にリロードします");
  });
  EventsOn("ON_FAIL_POLLING", (error) => {
    TonakoManager.getInstance.setErrorState(error.Error());
  });
  EventsOn("ON_START_BATTLE", async (tempArenaInfo: core.TempArenaInfo) => {
    TonakoManager.getInstance.setLoadingState("戦闘データを読み込み中");

    try {
      const battle =  await FetchBattle(tempArenaInfo)
      storedBattle.set(battle);
      TonakoManager.getInstance.setHidden();
    } catch (error) {
      TonakoManager.getInstance.setErrorState(error as string);
      return;
    }
  });

  const main = async () => {
    TonakoManager.getInstance.setLoadingState("設定ファイルの読み込み中");
    
    try {
      const displayPrefString = await LoadDisplayPref()

      let displayPref: DisplayPref;
      if (displayPrefString === "") {
        displayPref = DEFAULT_DISPLAY_PREF;
      } else {
        displayPref = JSON.parse(displayPrefString) as DisplayPref;
      }
      
      storedDisplayPref.set(displayPref);
    } catch (error) {
      TonakoManager.getInstance.setErrorState(error as string);
      return;
    }

    try {
      const gameClientPath = await LoadGameClientPath();
      storedGameClientPath.set(gameClientPath);
    } catch (error) {
      TonakoManager.getInstance.setErrorState(error as string);
      return;
    }

    TonakoManager.getInstance.setLoadingState("艦艇・マップの読み込み中");
    try {
      await Prefetch();
    } catch (error) {
      TonakoManager.getInstance.setErrorState(error as string);
      return;
    }

    StartPollingMatch();
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
        {#if page === "stats"}
          <StatsPage />
        {:else if page === "pref"}
          <PrefPage />
        {:else if page === "info"}
          <InfoPage />
        {/if}
      </div>
    </div>
  </div>
</main>
