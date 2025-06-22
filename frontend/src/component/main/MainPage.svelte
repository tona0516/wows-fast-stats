<script lang="ts">
  import { FetchProxy } from "src/lib/FetchProxy";
  import { Notifier } from "src/lib/Notifier";
  import {
    storedBattle,
    storedConfig,
    storedInstallPathError,
  } from "src/stores";
  import { LogInfo } from "wailsjs/go/main/App";

  import { RowPattern } from "src/lib/RowPattern";
  import { ColumnProvider } from "src/lib/column/ColumnProvider";

  let isLoading = false;

  $: categories = ColumnProvider.getAllColumns($storedConfig);
  $: [basicColumns, shipColumns, overallColumns] = categories;
  $: shipColumnCount = shipColumns.columnCount();
  $: statsColumnCount = shipColumnCount + overallColumns.columnCount();

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
  <div>
    {#if $storedBattle}
      {@const teams = $storedBattle.teams}

      <div class="flex justify-center">
        <div class="stats shadow">
          <div class="stat">
            <div class="stat-title">マップ</div>
            <div class="stat-value text-lg">{$storedBattle.meta.arena}</div>
          </div>
          <div class="stat">
            <div class="stat-title">種別</div>
            <div class="stat-value text-lg">{$storedBattle.meta.type}</div>
          </div>
          <div class="stat">
            <div class="stat-title">取得時刻</div>
            <div class="stat-value text-lg">
              {new Date($storedBattle.meta.unixtime * 1000).toLocaleString()}
            </div>
          </div>
        </div>
      </div>

      <div class="flex">
        <div class="overflow-x-auto w-screen pb-4">
          <table class="table text-nowrap">
            {#each teams as team}
              {#if team.players.length !== 0}
                <thead>
                  <tr>
                    {#each categories as category}
                      {#if category.columnCount() > 0}
                        <th
                          class="p-1 text-center"
                          colspan={category.columnCount()}
                        >
                          {category.dispName()}
                        </th>
                      {/if}
                    {/each}
                  </tr>
                  <tr>
                    {#each categories as category}
                      {#each category as column}
                        {#if column.shouldShow()}
                          <th class="p-1 text-center">{column.header}</th>
                        {/if}
                      {/each}
                    {/each}
                  </tr>
                </thead>
                <tbody>
                  {#each team.players as player}
                    {@const statsPattern = $storedConfig.stats_pattern}
                    {@const rowPattern = RowPattern.derive(
                      player,
                      statsPattern,
                      statsColumnCount,
                      shipColumnCount,
                    )}
                    <tr>
                      {#each basicColumns as column}
                        <td class="p-1">
                          <svelte:component
                            this={column.svelteComponent()}
                            {column}
                            {player}
                          />
                        </td>
                      {/each}

                      {#if rowPattern === RowPattern.NO_COLUMN}
                        <td class="p-1 text-center" colspan={statsColumnCount}
                        ></td>
                      {:else if rowPattern === RowPattern.PRIVATE}
                        <td class="p-1 text-center" colspan={statsColumnCount}
                          >PRIVATE</td
                        >
                      {:else if rowPattern === RowPattern.NO_STATS}
                        <td class="p-1 text-center" colspan={statsColumnCount}
                          >N/A</td
                        >
                      {:else if rowPattern === RowPattern.NO_SHIP_STATS}
                        <td class="p-1 text-center" colspan={shipColumnCount}
                          >N/A</td
                        >
                        {#each overallColumns as column}
                          {#if column.shouldShow()}
                            <td class="py-1">
                              <svelte:component
                                this={column.svelteComponent()}
                                {column}
                                {player}
                              />
                            </td>
                          {/if}
                        {/each}
                      {:else}
                        {#each shipColumns as column}
                          {#if column.shouldShow()}
                            <td class="p-1">
                              <svelte:component
                                this={column.svelteComponent()}
                                {column}
                                {player}
                              />
                            </td>
                          {/if}
                        {/each}
                        {#each overallColumns as column}
                          {#if column.shouldShow()}
                            <td class="p-1">
                              <svelte:component
                                this={column.svelteComponent()}
                                {column}
                                {player}
                              />
                            </td>
                          {/if}
                        {/each}
                      {/if}
                    </tr>
                  {/each}
                </tbody>
              {/if}
            {/each}
          </table>
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
  </div>

  {#if isLoading}
    <div
      class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-20"
    >
      <span class="loading loading-spinner"></span>
    </div>
  {/if}
</div>
