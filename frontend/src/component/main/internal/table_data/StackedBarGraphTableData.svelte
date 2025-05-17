<script lang="ts">
import type { StackedBarGraphParam } from "src/lib/column/StackedBarGraphParam";
import type { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { data } from "wailsjs/go/models";

export let column: AbstractStatsColumn<StackedBarGraphParam>;
export let player: data.Player;

$: param = column.displayValue(player);
</script>

<td class="td-graph">
  <table class="charts-css bar hide-data stacked">
    <thead>
      {#each param.items as _}
        <th scope="col" />
      {/each}
    </thead>
    <tbody>
      <tr>
        {#each param.items as item}
          {@const value = item.value.toFixed(param.digit)}
          <td style="--size: calc({value}/100); --color: {item.colorCode};"
            ><span class="data">{value}</span><span class="tooltip"
              >{item.label}<br />{value}%</span
            ></td
          >
        {/each}
      </tr>
    </tbody>
  </table>
</td>

<style>
  :global(.td-graph) {
    min-width: 5em;
    max-width: 5em;
  }
</style>
