<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
import { LogDebug } from "../wailsjs/runtime/runtime";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

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

function color(usingTierGroup: string, tierGroup: string): string {
  return usingTierGroup === tierGroup
    ? Const.TIER_S_COLORS[tierGroup]
    : Const.TIER_P_COLORS[tierGroup];
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
              <td
                style="--size: calc({value}/100); --color: {color(
                  convertToKey(player.ship_info.tier),
                  key
                )};"
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
