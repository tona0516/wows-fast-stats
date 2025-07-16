<script lang="ts">
  import { FetchProxy } from "src/lib/FetchProxy";
  import { storedBattle, storedInstallPathError } from "src/stores";
  import { LogInfo, ShowMessageDialog } from "wailsjs/go/main/App";

  import TeamCompareBarChart from "./internal/TeamCompareBarChart.svelte";
  import AllPlayerTable from "./internal/AllPlayerTable.svelte";
  import BattleMetaInfo from "./internal/BattleMetaInfo.svelte";

  let isLoading = false;

  export const fetchBattle = async () => {
    try {
      isLoading = true;

      const start = new Date().getTime();
      await FetchProxy.getBattle();
      const elapsed = (new Date().getTime() - start) / 1000;

      LogInfo("fetch success", { "duration(s)": elapsed.toFixed(1) });
    } catch (error) {
      if (error instanceof Error) {
        ShowMessageDialog(error.message);
      }
    } finally {
      isLoading = false;
    }
  };
</script>

<div>
  {#if $storedBattle}
    <div class="flex">
      <AllPlayerTable teams={$storedBattle.teams} />
    </div>

    <div>
      <BattleMetaInfo meta={$storedBattle.meta} />
    </div>

    <div>
      <TeamCompareBarChart battle={$storedBattle} />
    </div>
  {:else}
    <p>
      {#if $storedInstallPathError}
        設定画面から初期設定を行ってください。
      {:else}
        戦闘中ではありません。開始時に自動的にリロードします。
      {/if}
    </p>
  {/if}

  {#if isLoading}
    <div class="flex h-screen items-center justify-center">
      <span class="loading loading-ring loading-xl"></span>
    </div>
  {/if}
</div>
