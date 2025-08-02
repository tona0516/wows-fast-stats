<script lang="ts">
  import { FetchProxy } from "src/lib/FetchProxy";
  import { storedBattle, storedInstallPathError } from "src/stores";
  import { LogInfo, ShowMessageDialog } from "wailsjs/go/main/App";

  import TeamCompareBarChart from "./internal/TeamCompareBarChart.svelte";
  import BattleMetaInfo from "./internal/BattleMetaInfo.svelte";
  import MainStatsTable from "./internal/MainStatsTable.svelte";
  import type { data } from "wailsjs/go/models";

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

  // 各艦種のメタ情報をまとめる
  const shipTypes = [
    { type: "cv", caption: "空母" },
    { type: "bb", caption: "戦艦" },
    { type: "cl", caption: "巡洋艦" },
    { type: "dd", caption: "駆逐艦" },
    { type: "ss", caption: "潜水艦" },
  ];

  const filterByShipType = (type: string) => (player: data.Player) =>
    player.ship_info.type === type;
</script>

<div>
  {#if $storedBattle}
    <div class="flex">
      <MainStatsTable teams={$storedBattle.teams} />
    </div>

    <div class="pt-2 flex flex-col items-center">
      <BattleMetaInfo meta={$storedBattle.meta} />
    </div>

    <div class="pt-2 flex flex-col items-center overflow-x-auto">
      <div class="grid 2xl:grid-cols-2 place-content-center gap-4">
        <TeamCompareBarChart battle={$storedBattle} />

        {#each shipTypes as { type, caption }}
          <TeamCompareBarChart
            battle={$storedBattle}
            {caption}
            filterFunc={filterByShipType(type)}
          />
        {/each}
      </div>
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
