<script lang="ts">
  import { storedConfig } from "src/stores";

  import { RowPattern } from "src/lib/RowPattern";
  import { ColumnProvider } from "src/lib/column/ColumnProvider";
  import type { data } from "wailsjs/go/models";

  export let teams: data.Team[];

  $: categories = ColumnProvider.getAllColumns($storedConfig);
  $: [basicColumns, shipColumns, overallColumns] = categories;
  $: shipColumnCount = shipColumns.columnCount();
  $: statsColumnCount = shipColumnCount + overallColumns.columnCount();
</script>

<div class="overflow-x-auto w-screen pb-4">
  <table class="table text-nowrap">
    {#each teams as team}
      {#if team.players.length !== 0}
        <thead>
          <tr>
            {#each categories as category}
              {#if category.columnCount() > 0}
                <th
                  class="px-1 py-0 text-center"
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
                  <th class="px-1 py-0 text-center">{column.header}</th>
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
                <td class="px-1 py-0">
                  <svelte:component
                    this={column.svelteComponent()}
                    {column}
                    {player}
                  />
                </td>
              {/each}

              {#if rowPattern === RowPattern.NO_COLUMN}
                <td class="px-1 py-0 text-center" colspan={statsColumnCount}
                ></td>
              {:else if rowPattern === RowPattern.PRIVATE}
                <td class="px-1 py-0 text-center" colspan={statsColumnCount}
                  >PRIVATE</td
                >
              {:else if rowPattern === RowPattern.NO_STATS}
                <td class="px-1 py-0 text-center" colspan={statsColumnCount}
                  >N/A</td
                >
              {:else if rowPattern === RowPattern.NO_SHIP_STATS}
                <td class="px-1 py-0 text-center" colspan={shipColumnCount}
                  >N/A</td
                >
                {#each overallColumns as column}
                  {#if column.shouldShow()}
                    <td class="px-1 py-0">
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
                    <td class="px-1 py-0">
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
                    <td class="px-1 py-0">
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
