<script lang="ts">
  import { storedConfig, storedTeamThreatLevels } from "src/stores";

  import { RowPattern } from "src/lib/RowPattern";
  import { ColumnProvider } from "src/lib/column/ColumnProvider";
  import type { data } from "wailsjs/go/models";
  import ColspanTableData from "./table_data/ColspanTableData.svelte";

  export let teams: data.Team[];

  $: categories = ColumnProvider.getAllColumns($storedConfig);
  $: [basicColumns, shipColumns, overallColumns] = categories;
  $: shipColumnCount = shipColumns.columnCount();
  $: statsColumnCount = shipColumnCount + overallColumns.columnCount();
  $: allColumnCount = basicColumns.columnCount() + statsColumnCount;
</script>

<div class="overflow-x-auto w-screen pb-4">
  <table class="table text-nowrap">
    {#each teams as team, i}
      {#if team.players.length !== 0}
        <thead>
          {#if $storedConfig.display.overall.threat_level && $storedTeamThreatLevels && $storedTeamThreatLevels[i]}
            {@const teamThreatLevel = $storedTeamThreatLevels[i]}
            <tr>
              <th colspan={allColumnCount}>
                戦力評価値平均: <span class="text-lg"
                  >{teamThreatLevel.average.toFixed(0)}</span
                >
                [確度:
                <span class="text-lg"
                  >{teamThreatLevel.accuracy.toFixed(0)}</span
                >%] [介護指数:
                <span class="text-lg"
                  >{teamThreatLevel.dissociationDegree.toFixed(0)}</span
                >%]
              </th>
            </tr>
          {/if}

          <tr>
            {#each categories as category}
              {#if category.columnCount() > 0}
                <th class="p-1 text-center" colspan={category.columnCount()}>
                  {category.dispName()}
                </th>
              {/if}
            {/each}
          </tr>
          <tr>
            {#each categories as category}
              {#each category as column}
                {#if column.needsShow()}
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
                <svelte:component
                  this={column.getTableDataComponent()}
                  {column}
                  {player}
                />
              {/each}

              {#if rowPattern === RowPattern.NO_COLUMN}
                <ColspanTableData colspan={statsColumnCount} text="" />
              {:else if rowPattern === RowPattern.PRIVATE}
                <ColspanTableData colspan={statsColumnCount} text="PRIVATE" />
              {:else if rowPattern === RowPattern.NO_STATS}
                <ColspanTableData colspan={statsColumnCount} text="N/A" />
              {:else if rowPattern === RowPattern.NO_SHIP_STATS}
                <ColspanTableData colspan={shipColumnCount} text="N/A" />
                {#each overallColumns as column}
                  <svelte:component
                    this={column.getTableDataComponent()}
                    {column}
                    {player}
                  />
                {/each}
              {:else}
                {#each shipColumns as column}
                  <svelte:component
                    this={column.getTableDataComponent()}
                    {column}
                    {player}
                  />
                {/each}
                {#each overallColumns as column}
                  <svelte:component
                    this={column.getTableDataComponent()}
                    {column}
                    {player}
                  />
                {/each}
              {/if}
            </tr>
          {/each}
        </tbody>
      {/if}
    {/each}
  </table>
</div>
