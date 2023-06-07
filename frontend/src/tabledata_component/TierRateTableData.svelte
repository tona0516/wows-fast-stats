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

function toColor(ownTierGroupKey: string, tierGroupKey: string): string {
  return ownTierGroupKey === tierGroupKey
    ? Const.TIER_S_COLORS[tierGroupKey]
    : Const.TIER_P_COLORS[tierGroupKey];
}

function toLabel(tierGroupKey: string) {
  switch (tierGroupKey) {
    case "low":
      return "1~4";
    case "middle":
      return "5~7";
    case "high":
      return "8~★";
    default:
      return "";
  }
}

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
  tierGroup: { [key: string]: number }
): StackedBarGraphParam[] {
  if (tierGroup === undefined) {
    return [];
  }

  return Object.keys(tierGroup).map((key) => {
    const ownTierGroup = toTierGroupKey(player.ship_info.tier);
    return {
      label: toLabel(key),
      color: toColor(ownTierGroup, key),
      value: tierGroup[key],
    };
  });
}

const digit = Const.DIGITS[key];

$: tierGroup = values(player, statsPattern, statsCatetory, key) as {
  [key: string]: number;
};

$: params = derive(player, tierGroup);
</script>

<td class="td-graph">
  <StackedBarGraph params="{params}" digit="{digit}" />
</td>
