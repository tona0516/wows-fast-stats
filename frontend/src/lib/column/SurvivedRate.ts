import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import { CssClass, type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class SurvivedRate extends AbstractSingleColumn<CommonKey> {
  constructor(
    private userConfig: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super();
  }

  displayKey(): CommonKey {
    return "survived_rate";
  }

  minDisplayName(): string {
    return "生存率(勝|負)";
  }

  fullDisplayName(): string {
    return "生存率 (勝利|敗北)";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays[this.category].survived_rate;
  }

  tdClass(player: domain.Player): string {
    return CssClass.TD_MULTI;
  }

  displayValue(player: domain.Player): string {
    const digit = this.userConfig.custom_digit.survived_rate;
    const stats = toPlayerStats(player, this.userConfig.stats_pattern)[
      this.category
    ];
    return `${stats.win_survived_rate.toFixed(
      digit,
    )}% | ${stats.lose_survived_rate.toFixed(digit)}%`;
  }

  textColorCode(player: domain.Player): string {
    return "";
  }
}
