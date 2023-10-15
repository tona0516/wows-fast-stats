import { DispName } from "src/lib/DispName";
import type { StackedBarGraphParam } from "src/lib/column/StackedBarGraphParam";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { IGraphColumn } from "src/lib/column/intetface/IGraphColumn";
import type { OverallKey, TierGroup } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class UsingTierRate
  extends AbstractColumn<OverallKey>
  implements IGraphColumn
{
  constructor(
    svelteComponent: any,
    private config: domain.UserConfig,
  ) {
    super("using_tier_rate", "T割合", "ティア別プレイ割合", 1, svelteComponent);
  }

  shouldShowColumn(): boolean {
    return this.config.displays.overall.using_tier_rate;
  }

  getGraphParam(player: domain.Player): StackedBarGraphParam {
    const digit = this.config.custom_digit.using_tier_rate;
    const tierRateGroup = toPlayerStats(player, this.config.stats_pattern)
      .overall.using_tier_rate;
    const ownTierGroup = this.toTierGroup(player.ship_info.tier);
    const colors = this.config.custom_color.tier;

    const items = DispName.TIER_GROUPS.toArray().map((tier) => {
      const colorCode =
        tier.key === ownTierGroup
          ? colors.own[tier.key]
          : colors.other[tier.key];
      const rate = tierRateGroup[tier.key];

      return { label: tier.value, colorCode: colorCode, value: rate };
    });

    return { digit: digit, items: items };
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
