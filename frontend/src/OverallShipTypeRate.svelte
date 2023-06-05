<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
import type { StatsPattern } from "./StatsPattern";
import { values } from "./util";
import type { StatsCategory } from "./StatsCategory";

export let player: vo.Player;
export let displayPattern: DisplayPattern;
export let statsPattern: StatsPattern;
export let statsCatetory: StatsCategory;

const digit = Const.DIGITS["ship_type_rate"];

function color(usingShipType: string, shipType: string): string {
  return usingShipType === shipType
    ? Const.TYPE_S_COLORS[shipType]
    : Const.TYPE_P_COLORS[shipType];
}

$: shipTypeGroup = values(
  player,
  displayPattern,
  statsPattern,
  statsCatetory,
  "using_ship_type_rate"
);
</script>

{#if shipTypeGroup !== undefined}
  {@const keys = Object.keys(shipTypeGroup)}

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
            {@const value = shipTypeGroup[key].toFixed(digit)}
            <td
              style="--size: calc({value}/100); --color: {color(
                player.ship_info.type,
                key
              )};"
              ><span class="data">{value}</span><span class="tooltip"
                >{key.toUpperCase()}<br />{value}%</span
              ></td
            >
          {/each}
        </tr>
      </tbody>
    </table>
  </td>
{/if}
