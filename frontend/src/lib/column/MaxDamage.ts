import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import {
  CssClass,
  type CommonStatsKey,
  type StatsCategory,
} from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class MaxDamage extends AbstractSingleColumn<CommonStatsKey> {
  constructor(
    private userConfig: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super();
  }

  displayKey(): CommonStatsKey {
    return "max_damage";
  }

  minDisplayName(): string {
    return "最大Dmg";
  }

  fullDisplayName(): string {
    return "最大ダメージ";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays[this.category].max_damage;
  }

  tdClass(player: domain.Player): string {
    return CssClass.TD_NUM;
  }

  displayValue(player: domain.Player): string {
    const value = toPlayerStats(player, this.userConfig.stats_pattern)[
      this.category
    ].max_damage;
    const digit = this.userConfig.custom_digit.max_damage;
    return value.toFixed(digit);
  }

  textColorCode(player: domain.Player): string {
    return "";
  }
}
