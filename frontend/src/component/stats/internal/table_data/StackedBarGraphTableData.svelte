<script lang="ts">
  import type { StackedBarGraphParam } from "src/lib/column/StackedBarGraphParam";
  import type { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
  import type { data } from "wailsjs/go/models";

  export let column: AbstractStatsColumn<StackedBarGraphParam>;
  export let player: data.Player;

  $: param = column.displayValue(player);
</script>

{#if column.needsShow()}
  <td
    class="p-1"
    style="background-color: {column.getBackgroundColorCode(player) ?? ''}"
  >
    <div class="w-20">
      <table class="charts-css bar hide-data stacked">
        <thead>
          {#each param.items as _}
            <th scope="col" />
          {/each}
        </thead>
        <tbody>
          <tr>
            {#each param.items as item}
              {@const value = item.value.format(param.digit)}
              <td style="--size: calc({value}/100); --color: {item.colorCode};"
                ><span class="data">{value}</span><span class="tooltip"
                  >{item.label}<br />{value}%</span
                ></td
              >
            {/each}
          </tr>
        </tbody>
      </table>
    </div>
  </td>
{/if}
