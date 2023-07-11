<script lang="ts">
import type { StackedBarGraphParam } from "./StackedBarGraphParam";

export let param: StackedBarGraphParam;
</script>

<table class="charts-css bar hide-data stacked">
  <thead>
    {#each param.items as _}
      <th scope="col"></th>
    {/each}
  </thead>
  <tbody>
    <tr>
      {#each param.items as item}
        {@const value = item.value.toFixed(param.digit)}
        <td style="--size: calc({value}/100); --color: {item.color};"
          ><span class="data">{value}</span><span class="tooltip"
            >{item.label}<br />{value}%</span
          ></td
        >
      {/each}
    </tr>
  </tbody>
</table>
