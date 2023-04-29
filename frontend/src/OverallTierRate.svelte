<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;
</script>

<!-- using tier rate -->
{#if config.displays.overall.using_tier_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    {@const low = player.player_stats.using_tier_rate.low.toFixed(1)}
    {@const middle = player.player_stats.using_tier_rate.middle.toFixed(1)}
    {@const high = player.player_stats.using_tier_rate.high.toFixed(1)}

    <td class="td-graph">
      <table class="charts-css bar hide-data stacked">
        <thead>
          <th scope="col"> #1 </th>
          <th scope="col"> #2 </th>
          <th scope="col"> #3 </th>
        </thead>
        <tbody>
          <tr>
            <td style="--size: calc({low}/100); --color: #f9344c;"
              ><span class="data">{low}</span><span class="tooltip"
                >1~4<br />{low}%</span
              ></td
            >
            <td style="--size: calc({middle}/100); --color: #33a65e;"
              ><span class="data">{middle}</span><span class="tooltip"
                >5~7<br />{middle}%</span
              ></td
            >
            <td style="--size: calc({high}/100); --color: #1d86ae;"
              ><span class="data">{high}</span><span class="tooltip"
                >8~★<br />{high}%</span
              ></td
            >
          </tr>
        </tbody>
      </table>
    </td>
  {/if}
{/if}
