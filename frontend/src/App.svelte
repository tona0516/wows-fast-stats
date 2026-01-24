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
    LoadPref,
    StartPollingMatch,
    Prefetch,
  } from "@wails/go/main/App";
  import type { core } from "@wails/go/models";
  import { EventsOn } from "@wails/runtime/runtime";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import { TonakoManager } from "@libs/TonakoManager";
  import ShipDetailModal from "@components/modals/ShipDetailModal.svelte";
  import PlayerDetailModal from "@components/modals/PlayerDetailModal.svelte";

  let page: Page = "stats";

  $: {
    // @ts-ignore
    document.body.style.zoom = ($storedPref?.zoomRate || 1.0) / 100;
  }

  onMount(() => {
    themeChange(false);
  });

  EventsOn("ON_START_PREFETCH", (message: string) => {
    TonakoManager.getInstance.setLoadingState(message);
  });
  EventsOn("ON_START_PREFETCH_FAILURE", (message: string) => {
    TonakoManager.getInstance.setErrorState(message);
  });
  EventsOn("ON_PROMOTE", (message: string) => {
    TonakoManager.getInstance.setPromoteState(message);
  });
  EventsOn("ON_START_POLLING", (message: string) => {
    TonakoManager.getInstance.setStandbyState(message);
  });
  EventsOn("ON_START_BATTLE", (message: string) => {
    TonakoManager.getInstance.setLoadingState(message);
  });
  EventsOn("ON_FETCH_BATTLE_SUCCESS", (battle: core.Battle) => {
    TonakoManager.getInstance.setHidden();
    storedBattle.set(battle);
  });
  EventsOn("ON_FETCH_BATTLE_FAILURE", (message: string) => {
    TonakoManager.getInstance.setErrorState(message);
  });

  const main = async () => {
    try {
      const pref = await LoadPref();
      storedPref.set(pref);

      await Prefetch();
    } catch (error) {
      TonakoManager.getInstance.setErrorState(error as string);
      return
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
