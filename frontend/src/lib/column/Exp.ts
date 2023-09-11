import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import { CssClass, type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class Exp extends AbstractSingleColumn<CommonKey> {
  constructor(
    private userConfig: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super();
  }

  displayKey(): CommonKey {
    return "exp";
  }

  minDisplayName(): string {
    return "Exp";
  }

  fullDisplayName(): string {
    return "平均取得経験値(プレミアム補正含む)";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays[this.category].exp;
  }

  tdClass(player: domain.Player): string {
    return CssClass.TD_NUM;
  }

  displayValue(player: domain.Player): string {
    const digit = this.userConfig.custom_digit.exp;
    const value = toPlayerStats(player, this.userConfig.stats_pattern)[
      this.category
    ].exp;
    return value.toFixed(digit);
  }

  textColorCode(player: domain.Player): string {
    return "";
  }
}
