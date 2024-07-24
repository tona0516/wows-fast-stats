import StackedBarGraphTableData from "src/component/main/internal/table_data/StackedBarGraphTableData.svelte";
import { DispName } from "src/lib/DispName";
import type { StackedBarGraphParam } from "src/lib/column/StackedBarGraphParam";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { TierGroup } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class UsingTierRate extends AbstractStatsColumn<StackedBarGraphParam> {
  constructor(config: data.UserConfigV2) {
    super("using_tier_rate", 1, config, "overall");
  }

  displayValue(player: data.Player): StackedBarGraphParam {
    const tierRateGroup = this.playerStats(player).overall.using_tier_rate;
    const ownTierGroup = this.toTierGroup(player.ship_info.tier);
    const colors = this.config.color.tier;

    const items = DispName.TIER_GROUPS.toArray().map((it) => {
      const colorCode =
        it.key === ownTierGroup ? colors.own[it.key] : colors.other[it.key];
      const rate = tierRateGroup[it.key];

      return { label: it.value, colorCode: colorCode, value: rate };
    });

    return { digit: this.digit(), items: items };
  }

  svelteComponent() {
    return StackedBarGraphTableData;
  }

  private toTierGroup(tier: number): TierGroup | undefined {
    if (tier >= 1 && tier <= 4) {
      return "low";
    }
    if (tier >= 5 && tier <= 7) {
      return "middle";
    }
    if (tier >= 8) {
      return "high";
    }

    return undefined;
  }
}
