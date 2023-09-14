<script lang="ts">
  import { RowPattern } from "src/lib/types";
  import { tableColumns, toPlayerStats } from "src/lib/util";
  import type { domain } from "wailsjs/go/models";

  export let teams: domain.Team[];
  export let userConfig: domain.UserConfig;

  $: categories = tableColumns(userConfig);
  $: [basicColumns, shipColumns, overallColumns] = categories;
  $: shipColumnCount = shipColumns.columnCount();
  $: allColumnCount = shipColumnCount + overallColumns.columnCount();

  const decideRowPattern = (
    player: domain.Player,
    statsPattern: string,
    allColumnCount: number,
  ): RowPattern => {
    if (allColumnCount === 0) {
      return RowPattern.NO_COLUMN;
    }

    if (player.player_info.is_hidden === true) {
      return RowPattern.PRIVATE;
    }

    const stats = toPlayerStats(player, statsPattern);
    if (player.player_info.id === 0 || stats.overall.battles === 0) {
      return RowPattern.NO_DATA;
    }

    if (stats.ship.battles === 0) {
      return RowPattern.NO_SHIP_STATS;
    }

    return RowPattern.FULL;
  };
</script>

<div class="table-responsive">
  <table class="table table-sm table-text-color">
    {#each teams as team}
      <thead>
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
              {#if column.shouldShowColumn()}
                <th colspan={column.countInnerColumn()}
                  >{column.minDisplayName()}</th
                >
              {/if}
            {/each}
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each team.players as player}
          {@const statsPattern = userConfig.stats_pattern}
          {@const rowPattern = decideRowPattern(
            player,
            statsPattern,
            allColumnCount,
          )}
          <tr>
            {#each basicColumns as column}
              <svelte:component
                this={column.svelteComponent()}
                {column}
                {player}
                on:UpdateAlertPlayer
                on:RemoveAlertPlayer
                on:CheckPlayer
              />
            {/each}

            {#if rowPattern === RowPattern.PRIVATE}
              <td class="no_data omit" colspan={allColumnCount}>PRIVATE</td>
            {:else if rowPattern === RowPattern.NO_DATA}
              <td class="no_data omit" colspan={allColumnCount}>N/A</td>
            {:else if rowPattern === RowPattern.NO_SHIP_STATS}
              <td class="no_data omit" colspan={shipColumnCount}>N/A</td>
              {#each overallColumns as column}
                {#if column.shouldShowColumn()}
                  <svelte:component
                    this={column.svelteComponent()}
                    {column}
                    {player}
                  />
                {/if}
              {/each}
            {:else if rowPattern === RowPattern.FULL}
              {#each shipColumns as column}
                {#if column.shouldShowColumn()}
                  <svelte:component
                    this={column.svelteComponent()}
                    {column}
                    {player}
                  />
                {/if}
              {/each}

              {#each overallColumns as column}
                {#if column.shouldShowColumn()}
                  <svelte:component
                    this={column.svelteComponent()}
                    {column}
                    {player}
                  />
                {/if}
              {/each}
            {:else}
              <!-- Note: NO_COLUMN -->
            {/if}
          </tr>
        {/each}
      </tbody>
    {/each}
  </table>
</div>

<style>
  :global(.table-responsive) {
    max-width: 100%;
  }
</style>
