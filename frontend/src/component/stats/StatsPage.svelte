<script lang="ts">
  import { FetchProxy } from "src/lib/FetchProxy";
  import { storedBattle, storedInstallPathError } from "src/stores";
  import { LogInfo } from "wailsjs/go/main/App";

  import BattleMetaInfo from "./internal/BattleMetaInfo.svelte";
  import MainStatsTable from "./internal/MainStatsTable.svelte";
  import MessagingTonako from "./internal/MessagingTonako.svelte";
  import { Tonako } from "./internal/Tonako";

  let isLoading = false;
  let errorText = "";

  export const fetchBattle = async () => {
    try {
      isLoading = true;

      const start = new Date().getTime();
      await FetchProxy.getBattle();
      const elapsed = (new Date().getTime() - start) / 1000;

      LogInfo("fetch success", { "duration(s)": elapsed.toFixed(1) });
    } catch (error) {
      showError(error as string);
    } finally {
      isLoading = false;
    }
  };

  export const showError = (error: string) => {
    errorText = error;
  };
</script>

<div>
  {#if errorText.length > 0}
    <MessagingTonako tonako={Tonako.Sorry} message={errorText} />
  {:else if isLoading}
    <div class="flex w-full h-screen items-center justify-center">
      <span class="loading loading-ring loading-xl"></span>
    </div>
  {:else if $storedBattle}
    <div class="pt-2 flex flex-col items-center">
      <BattleMetaInfo meta={$storedBattle.meta} />
    </div>

    <div class="flex">
      <MainStatsTable teams={$storedBattle.teams} />
    </div>
  {:else if $storedInstallPathError}
    <MessagingTonako
      tonako={Tonako.Pointing}
      message="設定画面から初期設定をおこなってください"
    />
  {:else}
    <MessagingTonako
      tonako={Tonako.Standby}
      message="戦闘開始時に自動的にリロードします"
    />
  {/if}
</div>
