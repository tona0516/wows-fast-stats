<script lang="ts">
import type { vo } from "../../wailsjs/go/models";
import { Const } from "../Const";
import type { StackedBarGraphParam } from "../StackedBarGraphParam";
import StackedBarGraph from "../other_component/StackedBarGraph.svelte";
import type { StatsCategory } from "../enums";
import { values } from "../util";

export let player: vo.Player;
export let statsPattern: string;
export let statsCatetory: StatsCategory;
export let key: string;

const digit = Const.DIGITS[key];

function toColor(ownShipType: string, shipTypeKey: string): string {
  return ownShipType === shipTypeKey
    ? Const.TYPE_S_COLORS[shipTypeKey]
    : Const.TYPE_P_COLORS[shipTypeKey];
}

function derive(
  player: vo.Player,
  shipTypeGroup: { [key: string]: number }
): StackedBarGraphParam[] {
  if (shipTypeGroup === undefined) {
    return [];
  }

  return Object.keys(shipTypeGroup).map((key) => {
    return {
      label: key,
      color: toColor(player.ship_info.type, key),
      value: shipTypeGroup[key],
    };
  });
}

$: shipTypeGroup = values(player, statsPattern, statsCatetory, key) as {
  [key: string]: number;
};

$: params = derive(player, shipTypeGroup);
</script>

<td class="td-graph">
  <StackedBarGraph params="{params}" digit="{digit}" />
</td>
