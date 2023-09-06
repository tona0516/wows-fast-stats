import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import { CssClass, type OverallOnlyKey } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class AvgTier extends AbstractSingleColumn<OverallOnlyKey> {
  constructor(private userConfig: domain.UserConfig) {
    super();
  }

  displayKey(): OverallOnlyKey {
    return "avg_tier";
  }

  minDisplayName(): string {
    return "平均T";
  }

  fullDisplayName(): string {
    return "平均Tier";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays.overall.avg_tier;
  }

  tdClass(player: domain.Player): string {
    return CssClass.TD_NUM;
  }

  displayValue(player: domain.Player): string {
    const digit = this.userConfig.custom_digit.avg_tier;
    const value = toPlayerStats(player, this.userConfig.stats_pattern).overall
      .avg_tier;
    return value.toFixed(digit);
  }

  textColorCode(player: domain.Player): string {
    return "";
  }
}
