<script lang="ts">
import BattleMeta from "src/component/main/internal/BattleMeta.svelte";
import StatisticsTable from "src/component/main/internal/StatsTable.svelte";
import { FetchProxy } from "src/lib/FetchProxy";
import { Notifier } from "src/lib/Notifier";
import { storedBattle, storedConfig, storedInstallPathError } from "src/stores";
import { LogInfo } from "wailsjs/go/main/App";
import Menu from "./internal/Menu.svelte";

let menu: Menu | undefined;
let isLoading = false;

export const fetchBattle = async () => {
  try {
    isLoading = true;

    const start = new Date().getTime();
    await FetchProxy.getBattle();
    const elapsed = (new Date().getTime() - start) / 1000;

    Notifier.success(`データ取得完了: ${elapsed.toFixed(1)}秒`);

    LogInfo("fetch success", { "duration(s)": elapsed.toFixed(1) });
  } catch (error) {
    Notifier.failure(error);
  } finally {
    isLoading = false;
  }
};
</script>

<div>
  <div class="flex">
    <Menu bind:this={menu} />
  </div>

  <div>
    {#if $storedBattle}
      {@const teams = $storedBattle.teams}
      {@const meta = $storedBattle.meta}
      {@const config = $storedConfig}

      <div class="flex">
        <StatisticsTable {teams} {config} />
      </div>

      <div class="flex">
        <BattleMeta {meta} />
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
  </div>

  {#if isLoading}
    <div
      class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-20"
    >
      <span class="loading loading-spinner"></span>
    </div>
  {/if}
</div>
