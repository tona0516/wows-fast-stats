<script lang="ts">
  import { RowPattern } from "src/lib/RowPattern";
  import { ColumnProvider } from "src/lib/column/ColumnProvider";
  import type { data } from "wailsjs/go/models";

  export let teams: data.Team[];
  export let config: data.UserConfigV2;

  $: categories = ColumnProvider.getAllColumns(config);
  $: [basicColumns, shipColumns, overallColumns] = categories;
  $: shipColumnCount = shipColumns.columnCount();
  $: statsColumnCount = shipColumnCount + overallColumns.columnCount();
</script>

<div class="overflow-x-auto w-screen py-4">
  <table class="table text-nowrap">
    {#each teams as team}
      {#if team.players.length !== 0}
        <thead>
          <tr>
            {#each categories as category}
              {#if category.columnCount() > 0}
                <th class="p-1 text-center" colspan={category.columnCount()}
                  >{category.dispName()}</th
                >
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
            {@const statsPattern = config.stats_pattern}
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
                <td class="p-1 bg-base-300 text-center" colspan={statsColumnCount}></td>
              {:else if rowPattern === RowPattern.PRIVATE}
                <td class="p-1 bg-base-300 text-center" colspan={statsColumnCount}
                  >PRIVATE</td
                >
              {:else if rowPattern === RowPattern.NO_STATS}
                <td class="p-1 bg-base-300 text-center" colspan={statsColumnCount}>N/A</td>
              {:else if rowPattern === RowPattern.NO_SHIP_STATS}
                <td class="p-1 bg-base-300 text-center" colspan={shipColumnCount}>N/A</td>
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
