<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

let digit = Const.DIGITS["ship_type_rate"];

function color(usingShipType: string, shipType: string): string {
  return usingShipType === shipType
    ? Const.TYPE_S_COLORS[shipType]
    : Const.TYPE_P_COLORS[shipType];
}
</script>

{#if config.displays.overall.using_ship_type_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    {@const keys = Object.keys(player.overall_stats.using_ship_type_rate)}

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
                player.overall_stats.using_ship_type_rate[key].toFixed(digit)}
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
{/if}
