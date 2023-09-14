import { DispName } from "src/lib/DispName";
import { AbstractGraphColumn } from "src/lib/column/intetface/AbstractGraphColumn";
import type { OverallKey, TierGroup } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type {
  StackedBarGraphItem,
  StackedBarGraphParam,
} from "src/lib/value_object/StackedBarGraphParam";
import type { domain } from "wailsjs/go/models";

export class UsingTierRate extends AbstractGraphColumn<OverallKey> {
  constructor(private userConfig: domain.UserConfig) {
    super();
  }

  displayKey(): OverallKey {
    return "using_tier_rate";
  }

  minDisplayName(): string {
    return "T割合";
  }

  fullDisplayName(): string {
    return "ティア別プレイ割合";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays.overall.using_tier_rate;
  }

  displayValue(player: domain.Player): StackedBarGraphParam {
    const digit = this.userConfig.custom_digit.using_tier_rate;
    const tierRateGroup = toPlayerStats(player, this.userConfig.stats_pattern)
      .overall.using_tier_rate;
    const ownTierGroup = this.toTierGroup(player.ship_info.tier);
    const colors = this.userConfig.custom_color.tier;

    let items: StackedBarGraphItem[] = [];
    DispName.TIER_GROUPS.forEach((value, key) => {
      const rate = tierRateGroup[key];
      const colorCode =
        key === ownTierGroup ? colors.own[key] : colors.other[key];
      items.push({ label: value, colorCode: colorCode, value: rate });
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
