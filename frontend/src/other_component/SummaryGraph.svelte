<script lang="ts">
  import type { Summary } from "src/lib/Summary";
  import { storedUserConfig } from "src/stores";

  export let summary: Summary;
</script>

<table
  id="summary-table"
  class="charts-css column multiple show-labels data-spacing-5 show-primary-axis"
>
  <thead>
    <tr>
      {#each Object.keys(summary.friends) as label}
        <th scope="col">{label}</th>
      {/each}
    </tr>
  </thead>
  <tbody>
    {#each Object.keys(summary.friends) as label}
      {@const friend = summary.friends[label]}
      {@const enemy = summary.enemies[label]}
      {@const max = Math.max(friend, enemy)}
      <tr>
        <th scope="row">{label}</th>
        <td style="--size: {friend / max};">{friend}</td>
        <td style="--size: {enemy / max};">{enemy}</td>
      </tr>
    {/each}
  </tbody>
</table>
<ul
  id="summary-legend"
  class="charts-css legend legend-inline legend-square center"
  style="font-size: {$storedUserConfig.font_size};"
>
  <li>味方チーム平均</li>
  <li>敵チーム平均</li>
</ul>

<style>
  #summary-table {
    --color-1: #518517;
    --color-2: #a41200;
    height: 120px;
    max-width: 800px;
  }

  #summary-legend {
    --color-1: #518517;
    --color-2: #a41200;
    max-width: 400px;
    padding: 0;
  }
</style>
