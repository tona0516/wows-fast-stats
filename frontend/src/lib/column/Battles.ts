import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import {
  CssClass,
  type CommonStatsKey,
  type StatsCategory,
} from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class Battles extends AbstractSingleColumn<CommonStatsKey> {
  constructor(
    private userConfig: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super();
  }

  displayKey(): CommonStatsKey {
    return "battles";
  }

  minDisplayName(): string {
    return "戦闘数";
  }

  fullDisplayName(): string {
    return "戦闘数";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays[this.category].battles;
  }

  tdClass(player: domain.Player): string {
    return CssClass.TD_NUM;
  }

  displayValue(player: domain.Player): string {
    const digit = this.userConfig.custom_digit.battles;
    const value = toPlayerStats(player, this.userConfig.stats_pattern)[
      this.category
    ].battles;
    return value.toFixed(digit);
  }

  textColorCode(player: domain.Player): string {
    return "";
  }
}
