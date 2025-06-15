<script lang="ts">
import { RowPattern } from "src/lib/RowPattern";
import { ColumnProvider } from "src/lib/column/ColumnProvider";
import { storedTeamThreatLevels } from "src/stores";
import type { data } from "wailsjs/go/models";

export let teams: data.Team[];
export let config: data.UserConfigV2;

$: categories = ColumnProvider.getAllColumns(config);
$: [basicColumns, shipColumns, overallColumns] = categories;
$: shipColumnCount = shipColumns.columnCount();
$: statsColumnCount = shipColumnCount + overallColumns.columnCount();
$: allColumnCount = basicColumns.columnCount() + statsColumnCount;
</script>

<div class="overflow-x-auto w-screen">
  <table class="table text-nowrap">
    {#each teams as team, i}
      {#if team.players.length !== 0}
        <thead>
          {#if config.display.overall.threat_level && $storedTeamThreatLevels && $storedTeamThreatLevels[i]}
            {@const teamThreatLevel = $storedTeamThreatLevels[i]}
            <tr>
              <th colspan={allColumnCount}>
                戦力評価値平均 : {teamThreatLevel.average.toFixed(0)}
                [確度 : {teamThreatLevel.accuracy.toFixed(0)}%] [介護指数 : {teamThreatLevel.dissociationDegree.toFixed(
                  0,
                )}%]
              </th>
            </tr>
          {/if}
          <tr>
            {#each categories as category}
              {#if category.columnCount() > 0}
                <th colspan={category.columnCount()}>{category.dispName()}</th>
              {/if}
            {/each}
          </tr>
          <tr>
            {#each categories as category}
              {#each category as column}
                {#if column.shouldShow()}
                  <th>{column.header}</th>
                {/if}
              {/each}
            {/each}
          </tr>
        </thead>

        <tbody>
          {#each team.players as player}
            {@const statsPattern = config.stats_pattern}
            {@const rowPattern = RowPattern.derive(
              player,
              statsPattern,
              statsColumnCount,
              shipColumnCount,
            )}
            <tr>
              {#each basicColumns as column}
                <svelte:component
                  this={column.svelteComponent()}
                  {column}
                  {player}
                  on:EditAlertPlayer
                  on:RemoveAlertPlayer
                />
              {/each}

              {#if rowPattern === RowPattern.NO_COLUMN}
                <td class="text-center" colspan={statsColumnCount}></td>
              {:else if rowPattern === RowPattern.PRIVATE}
                <td class="text-center" colspan={statsColumnCount}>PRIVATE</td>
              {:else if rowPattern === RowPattern.NO_STATS}
                <td class="text-center" colspan={statsColumnCount}>N/A</td>
              {:else if rowPattern === RowPattern.NO_SHIP_STATS}
                <td class="text-center" colspan={shipColumnCount}>N/A</td>
                {#each overallColumns as column}
                  {#if column.shouldShow()}
                    <svelte:component
                      this={column.svelteComponent()}
                      {column}
                      {player}
                    />
                  {/if}
                {/each}
              {:else}
                {#each shipColumns as column}
                  {#if column.shouldShow()}
                    <svelte:component
                      this={column.svelteComponent()}
                      {column}
                      {player}
                    />
                  {/if}
                {/each}

                {#each overallColumns as column}
                  {#if column.shouldShow()}
                    <svelte:component
                      this={column.svelteComponent()}
                      {column}
                      {player}
                    />
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
