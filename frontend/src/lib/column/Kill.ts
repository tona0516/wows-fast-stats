import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import {
  CssClass,
  type CommonStatsKey,
  type StatsCategory,
} from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class Kill extends AbstractSingleColumn<CommonStatsKey> {
  constructor(
    private userConfig: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super();
  }

  displayKey(): CommonStatsKey {
    return "kill";
  }

  minDisplayName(): string {
    return "撃沈";
  }

  fullDisplayName(): string {
    return "平均撃沈数";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays[this.category].kill;
  }

  tdClass(player: domain.Player): string {
    return CssClass.TD_NUM;
  }

  displayValue(player: domain.Player): string {
    const digit = this.userConfig.custom_digit.kill;
    const value = toPlayerStats(player, this.userConfig.stats_pattern)[
      this.category
    ].kill;
    return value.toFixed(digit);
  }

  textColorCode(player: domain.Player): string {
    return "";
  }
}
