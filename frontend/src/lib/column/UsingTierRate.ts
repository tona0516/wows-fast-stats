import { DispName } from "src/lib/DispName";
import type {
  StackedBarGraphItem,
  StackedBarGraphParam,
} from "src/lib/StackedBarGraphParam";
import { AbstractGraphColumn } from "src/lib/column/intetface/AbstractGraphColumn";
import type { OverallOnlyKey } from "src/lib/types";
import { tierString, toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class UsingTierRate extends AbstractGraphColumn<OverallOnlyKey> {
  constructor(private userConfig: domain.UserConfig) {
    super();
  }

  displayKey(): OverallOnlyKey {
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
    const ownTierGroup = tierString(player.ship_info.tier);
    const colors = this.userConfig.custom_color.tier;

    let items: StackedBarGraphItem[] = [];
    DispName.TIER_GROUPS.forEach((pair) => {
      const value = tierRateGroup[pair.first];
      const colorCode =
        pair.first === ownTierGroup
          ? colors.own[pair.first]
          : colors.other[pair.first];
      items.push({ label: pair.second, colorCode: colorCode, value: value });
    });

    return { digit: digit, items: items };
  }
}
