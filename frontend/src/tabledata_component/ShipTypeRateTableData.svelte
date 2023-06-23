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
export let customColor: vo.CustomColor;
export let customDigit: vo.CustomDigit;

function derive(
  player: vo.Player,
  shipTypeGroup: { [key: string]: number },
  customColor: vo.CustomColor
): StackedBarGraphParam[] {
  if (shipTypeGroup === undefined) {
    return [];
  }

  return Object.keys(shipTypeGroup).map((key) => {
    return {
      label: Const.SHIP_TYPE_LABELS[key] ?? "",
      color:
        key === player.ship_info.type
          ? customColor.ship_type.own[key]
          : customColor.ship_type.other[key],
      value: shipTypeGroup[key],
    };
  });
}

$: shipTypeGroup = values(player, statsPattern, statsCatetory, key) as {
  [key: string]: number;
};

$: params = derive(player, shipTypeGroup, customColor);

$: digit = customDigit[key];
</script>

<td class="td-graph">
  <StackedBarGraph params="{params}" digit="{digit}" />
</td>
