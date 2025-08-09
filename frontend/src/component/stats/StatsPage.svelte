<script lang="ts">
  import { storedBattle, storedInstallPathError } from "src/stores";
  import { Battle } from "wailsjs/go/main/App";

  import BattleMetaInfo from "./internal/BattleMetaInfo.svelte";
  import MainStatsTable from "./internal/MainStatsTable.svelte";
  import MessagingTonako from "./internal/MessagingTonako.svelte";
  import { Tonako } from "./internal/Tonako";

  let isLoading = false;
  let errorText = "";

  export const fetchBattle = async () => {
    try {
      isLoading = true;

      // Note: 過去のデータが影響してか値が0になってしまうためクリーンする
      storedBattle.set(undefined);

      // const cache = localStorage.getItem("cache");
      // TODO: あとで消す
      // if (cache) {
      //   const cachedBattle = JSON.parse(cache) as data.Battle;
      //   storedBattle.set(cachedBattle);
      //   return cachedBattle;
      // }

      const ret = await Battle();
      storedBattle.set(ret);

      localStorage.setItem("cache", JSON.stringify(ret));
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
