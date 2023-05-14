<script lang="ts">
import type { vo } from "wailsjs/go/models";
import type { DisplayPattern } from "./DisplayPattern";
import Const from "./Const";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

const usingShipColors: { [key: string]: string } = {
  ss: "#233B8B",
  dd: "#CCB914",
  cl: "#27853F",
  bb: "#CA1028",
  cv: "#5E2883",
};

const otherColors: { [key: string]: string } = {
  ss: "#B3CDE3",
  dd: "#FEE6AA",
  cl: "#CCEBC5",
  bb: "#FBB4C4",
  cv: "#CAB2D6",
};

function convertToKey(shipType: string): string {
  switch (shipType) {
    case "AirCarrier":
      return "cv";
    case "Battleship":
      return "bb";
    case "Cruiser":
      return "cl";
    case "Destroyer":
      return "dd";
    case "Submarine":
      return "ss";
    default:
      return "";
  }
}

let digit = Const.DIGITS["ship_type_rate"];
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
              {#if key === convertToKey(player.ship_info.type)}
                {@const color = usingShipColors[key]}
                <td style="--size: calc({value}/100); --color: {color};"
                  ><span class="data">{value}</span><span class="tooltip"
                    >{key.toUpperCase()}<br />{value}%</span
                  ></td
                >
              {:else}
                {@const color = otherColors[key]}
                <td style="--size: calc({value}/100); --color: {color};"
                  ><span class="data">{value}</span><span class="tooltip"
                    >{key.toUpperCase()}<br />{value}%</span
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
