<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

{#if config.displays.overall.using_ship_type_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    {@const ss = player.player_stats.using_ship_type_rate.ss.toFixed(1)}
    {@const dd = player.player_stats.using_ship_type_rate.dd.toFixed(1)}
    {@const cl = player.player_stats.using_ship_type_rate.cl.toFixed(1)}
    {@const bb = player.player_stats.using_ship_type_rate.bb.toFixed(1)}
    {@const cv = player.player_stats.using_ship_type_rate.cv.toFixed(1)}

    <td class="td-graph">
      <table class="charts-css bar hide-data stacked">
        <thead>
          <th scope="col"> #1 </th>
          <th scope="col"> #2 </th>
          <th scope="col"> #3 </th>
          <th scope="col"> #4 </th>
          <th scope="col"> #5 </th>
        </thead>
        <tbody>
          <tr>
            <td style="--size: calc({ss}/100); --color: #386cb0;"
              ><span class="data">{ss}</span><span class="tooltip"
                >SS<br />{ss}%</span
              ></td
            >
            <td style="--size: calc({dd}/100); --color: #fff231;"
              ><span class="data">{dd}</span><span class="tooltip"
                >DD<br />{dd}%</span
              ></td
            >
            <td style="--size: calc({cl}/100); --color: #99d02b;"
              ><span class="data">{cl}</span><span class="tooltip"
                >CL<br />{cl}%</span
              ></td
            >
            <td style="--size: calc({bb}/100); --color: #ff9914;"
              ><span class="data">{bb}</span><span class="tooltip"
                >BB<br />{bb}%</span
              ></td
            >
            <td style="--size: calc({cv}/100); --color: #a45aaa;"
              ><span class="data">{cv}</span><span class="tooltip"
                >CV<br />{cv}%</span
              ></td
            >
          </tr>
        </tbody>
      </table>
    </td>
  {/if}
{/if}
