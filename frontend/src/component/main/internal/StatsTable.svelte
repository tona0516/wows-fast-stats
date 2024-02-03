<script lang="ts">
  import UkTable from "src/component/common/uikit/UkTable.svelte";
  import { RowPattern } from "src/lib/RowPattern";
  import type { model } from "wailsjs/go/models";
  import { ColumnProvider } from "src/lib/column/ColumnProvider";
  import { FetchProxy } from "src/lib/FetchProxy";

  export let teams: model.Team[];
  export let config: model.UserConfig;

  $: categories = ColumnProvider.getAllColumns(config);
  $: [basicColumns, shipColumns, overallColumns] = categories;
  $: shipColumnCount = shipColumns.columnCount();
  $: allColumnCount = shipColumnCount + overallColumns.columnCount();
</script>

<UkTable>
  {#each teams as team}
    <thead>
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
        {@const statsPattern = config.stats_pattern}
        {@const rowPattern = RowPattern.derive(
          player,
          statsPattern,
          allColumnCount,
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
              on:CheckPlayer={() => FetchProxy.getExcludedPlayers()}
            />
          {/each}

          {#if rowPattern === RowPattern.NO_COLUMN}
            <td class="no_data" colspan={allColumnCount}></td>
          {:else if rowPattern === RowPattern.PRIVATE}
            <td class="no_data" colspan={allColumnCount}>PRIVATE</td>
          {:else if rowPattern === RowPattern.NO_STATS}
            <td class="no_data" colspan={allColumnCount}>N/A</td>
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
  {/each}
</UkTable>

<style>
  :global(.no_data) {
    max-width: 0px;
    text-align: center;
  }
</style>
