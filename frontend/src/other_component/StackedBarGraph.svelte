<script lang="ts">
import type { StackedBarGraphParam } from "../StackedBarGraphParam";

export let params: StackedBarGraphParam[];
export let digit: number;
</script>

<table class="charts-css bar hide-data stacked">
  <thead>
    {#each params as _}
      <th scope="col"></th>
    {/each}
  </thead>
  <tbody>
    <tr>
      {#each params as param}
        {@const value = param.value.toFixed(digit)}
        <td style="--size: calc({value}/100); --color: {param.color};"
          ><span class="data">{value}</span><span class="tooltip"
            >{param.label}<br />{value}%</span
          ></td
        >
      {/each}
    </tr>
  </tbody>
</table>
