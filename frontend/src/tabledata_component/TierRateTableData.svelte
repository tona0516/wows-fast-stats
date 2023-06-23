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

function toTierGroupKey(tier: number): string {
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

function derive(
  player: vo.Player,
  tierGroup: { [key: string]: number },
  customColor: vo.CustomColor
): StackedBarGraphParam[] {
  if (tierGroup === undefined) {
    return [];
  }

  return Object.keys(tierGroup).map((key) => {
    const ownTierGroup = toTierGroupKey(player.ship_info.tier);
    return {
      label: Const.TIER_GROUP_LABELS[key] ?? "",
      color:
        key === ownTierGroup
          ? customColor.tier.own[key]
          : customColor.tier.other[key],
      value: tierGroup[key],
    };
  });
}

$: tierGroup = values(player, statsPattern, statsCatetory, key) as {
  [key: string]: number;
};

$: params = derive(player, tierGroup, customColor);

$: digit = customDigit[key];
</script>

<td class="td-graph">
  <StackedBarGraph params="{params}" digit="{digit}" />
</td>
