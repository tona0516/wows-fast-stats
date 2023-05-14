<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

const usingColors: { [key: string]: string } = {
  low: "#A41200",
  middle: "#518517",
  high: "#04436D",
};

const otherColors: { [key: string]: string } = {
  low: "#FDCDB7",
  middle: "#E6F5B0",
  high: "#B3D7DD",
};

function convertToKey(tier: number): string {
  if (tier >= 1 && tier <= 4) {
    return "low";
  }
  if (tier >= 5 && tier <= 7) {
    return "middle";
  }
  if (tier >= 8) {
    return "high";
  }
  return "";
}

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
              {#if key === convertToKey(player.ship_info.tier)}
                {@const color = usingColors[key]}
                <td style="--size: calc({value}/100); --color: {color};"
                  ><span class="data">{value}</span><span class="tooltip"
                    >{texts[key]}<br />{value}%</span
                  ></td
                >
              {:else}
                {@const color = otherColors[key]}
                <td style="--size: calc({value}/100); --color: {color};"
                  ><span class="data">{value}</span><span class="tooltip"
                    >{texts[key]}<br />{value}%</span
                  ></td
                >
              {/if}
            {/each}
          </tr>
        </tbody>
      </table>
    </td>
  {/if}
{/if}
