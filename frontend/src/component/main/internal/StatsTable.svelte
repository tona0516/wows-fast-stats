<script lang="ts">
  import UkTable from "src/component/common/uikit/UkTable.svelte";
  import { RowPattern } from "src/lib/RowPattern";
  import type { data } from "wailsjs/go/models";
  import { ColumnProvider } from "src/lib/column/ColumnProvider";
  import { FetchProxy } from "src/lib/FetchProxy";
  import { storedTeamThreatLevels } from "src/stores";
  import type { StatsTableOptions } from "src/lib/StatsTableOptions";
  import StatsPatternSwitch from "src/component/main/internal/StatsPatternSwitch.svelte";

  export let teams: data.Team[];
  export let config: data.UserConfigV2;
  export let options: StatsTableOptions;

  $: categories = ColumnProvider.getAllColumns(config);
  $: [basicColumns, shipColumns, overallColumns] = categories;
  $: shipColumnCount = shipColumns.columnCount();
  $: statsColumnCount = shipColumnCount + overallColumns.columnCount();
  $: allColumnCount = basicColumns.columnCount() + statsColumnCount;
</script>

<div>
  <div class="uk-flex uk-flex-center">
    <StatsPatternSwitch {options} />
  </div>

  <div class="uk-flex uk-flex-center">
    <UkTable>
      {#each teams as team, i}
        {#if team.players.length !== 0}
          <thead>
            {#if config.display.overall.threat_level && $storedTeamThreatLevels && $storedTeamThreatLevels[i]}
              {@const teamThreatLevel = $storedTeamThreatLevels[i]}
              <tr>
                <th class="uk-text-center" colspan={allColumnCount}>
                  戦力評価値平均 : <span class="uk-text-large uk-text-bold"
                    >{teamThreatLevel.average.toFixed(0)}</span
                  >
                  [確度 :
                  <span class="uk-text-default uk-text-bold"
                    >{teamThreatLevel.accuracy.toFixed(0)}</span
                  >%] [介護指数 :
                  <span class="uk-text-default uk-text-bold"
                    >{teamThreatLevel.dissociationDegree.toFixed(0)}</span
                  >%]
                </th>
              </tr>
            {/if}
            <tr>
              {#each categories as category}
                {#if category.columnCount() > 0}
                  <th class="uk-text-center" colspan={category.columnCount()}
                    >{category.dispName()}</th
                  >
                {/if}
              {/each}
            </tr>
            <tr>
              {#each categories as category}
                {#each category as column}
                  {#if column.shouldShow()}
                    <th class="uk-text-center" colspan={column.innerColumnCount}
                      >{column.header}</th
                    >
                  {/if}
                {/each}
              {/each}
            </tr>
          </thead>

          <tbody>
            {#each team.players as player}
              {@const rowPattern = RowPattern.derive(
                player,
                options.statsPattern,
                statsColumnCount,
                shipColumnCount,
              )}
              <tr>
                {#each basicColumns as column}
                  <svelte:component
                    this={column.svelteComponent()}
                    {column}
                    {player}
                    {options}
                    on:EditAlertPlayer
                    on:RemoveAlertPlayer
                    on:CheckPlayer={() => FetchProxy.getExcludedPlayers()}
                  />
                {/each}

                {#if rowPattern === RowPattern.NO_COLUMN}
                  <td class="no_data" colspan={statsColumnCount}></td>
                {:else if rowPattern === RowPattern.PRIVATE}
                  <td class="no_data" colspan={statsColumnCount}>PRIVATE</td>
                {:else if rowPattern === RowPattern.NO_STATS}
                  <td class="no_data" colspan={statsColumnCount}>N/A</td>
                {:else if rowPattern === RowPattern.NO_SHIP_STATS}
                  <td class="no_data" colspan={shipColumnCount}>N/A</td>
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
    </UkTable>
  </div>
</div>

<style>
  :global(.no_data) {
    max-width: 0px;
    text-align: center;
  }
</style>
