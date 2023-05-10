<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

const colors: { [key: string]: string } = {
  low: "#f9344c",
  middle: "#33a65e",
  high: "#1d86ae",
};

const texts: { [key: string]: string } = {
  low: "1~4",
  middle: "5~7",
  high: "8~★",
};

let digit = Const.DIGITS["tier_rate"];
</script>

<!-- using tier rate -->
{#if config.displays.overall.using_tier_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    {@const keys = Object.keys(player.overall_stats.using_tier_rate)}

    <td class="td-graph">
      <table class="charts-css bar hide-data stacked">
        <thead>
          {#each keys as _}
            <th scope="col"></th>
          {/each}
        </thead>
        <tbody>
          <tr>
            {#each keys as key}
              {@const value =
                player.overall_stats.using_tier_rate[key].toFixed(digit)}
              <td style="--size: calc({value}/100); --color: {colors[key]};"
                ><span class="data">{value}</span><span class="tooltip"
                  >{texts[key]}<br />{value}%</span
                ></td
              >
            {/each}
          </tr>
        </tbody>
      </table>
    </td>
  {/if}
{/if}
