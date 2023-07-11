<script lang="ts">
import type { vo } from "../../wailsjs/go/models";
import { Const } from "../Const";
import type { StackedBarGraphParam } from "../other_component/stacked_bar/StackedBarGraphParam";
import StackedBarGraph from "../other_component/stacked_bar/StackedBarGraph.svelte";
import type { StatsCategory } from "../enums";
import { toTierGroup, values } from "../util";

export let player: vo.Player;
export let statsPattern: string;
export let statsCatetory: StatsCategory;
export let columnKey: string;
export let customColor: vo.CustomColor;
export let customDigit: vo.CustomDigit;

function toParam(
  player: vo.Player,
  tierGroup: { [key: string]: number },
  customColor: vo.CustomColor,
  customDigit: vo.CustomDigit
): StackedBarGraphParam {
  const digit = customDigit[columnKey];

  if (tierGroup === undefined) {
    return { digit: digit, items: [] };
  }

  const items = Object.keys(tierGroup).map((key) => {
    const ownTierGroup = toTierGroup(player.ship_info.tier);
    const label = Const.TIER_GROUP_LABELS[key] ?? "";
    const color =
      key === ownTierGroup
        ? customColor.tier.own[key]
        : customColor.tier.other[key];
    const value = tierGroup[key];

    return {
      label: label,
      color: color,
      value: value,
    };
  });

  return {
    digit: digit,
    items: items,
  };
}

$: tierGroup = values(player, statsPattern, statsCatetory, columnKey) as {
  [key: string]: number;
};

$: param = toParam(player, tierGroup, customColor, customDigit);
</script>

<td class="td-graph">
  <StackedBarGraph param="{param}" />
</td>
