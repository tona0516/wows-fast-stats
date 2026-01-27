<script lang="ts">
  import type { AbstractStatsColumn } from "@libs/columns/AbstractStatsColumn";
    import { storedDisplayPref } from "@libs/stores";
  import type { StackedBarChartParam } from "@libs/types";
  import type { core } from "@wails/go/models";

  export let column: AbstractStatsColumn<StackedBarChartParam[]>;
  export let player: core.Player;

  $: params = column.getDisplayValue(player);
</script>

{#if column.needsShow()}
  <td
    class="{$storedDisplayPref.showBoarder ? "border border-gray-500" : ""} p-1"
    style="background-color: {column.getBgColorCode(player)?.raw}"
  >
    <div class="w-20">
      <table class="charts-css bar hide-data stacked">
        <thead>
          {#each params as _}
            <th scope="col" />
          {/each}
        </thead>
        <tbody>
          <tr>
            {#each params as param}
              {@const value = param.value.toFixed(column.getDigit())}
              <td
                style="--size: calc({value}/100); --color: {param.colorCode
                  ?.raw};"
                ><span class="data">{value}</span><span class="tooltip"
                  >{param.label}<br />{value}%</span
                ></td
              >
            {/each}
          </tr>
        </tbody>
      </table>
    </div>
  </td>
{/if}
