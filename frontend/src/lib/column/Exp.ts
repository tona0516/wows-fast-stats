import { CssClass } from "src/lib/CssClass";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import { type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class Exp extends AbstractColumn<CommonKey> implements ISingleColumn {
  constructor(
    svelteComponent: any,
    private config: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super(
      "exp",
      "Exp",
      "平均取得経験値(プレミアム補正含む)",
      1,
      svelteComponent,
    );
  }

  shouldShowColumn(): boolean {
    return this.config.displays[this.category].exp;
  }

  getTdClass(_: domain.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: domain.Player): string {
    const digit = this.config.custom_digit.exp;
    const value = toPlayerStats(player, this.config.stats_pattern)[
      this.category
    ].exp;
    return value.toFixed(digit);
  }

  getTextColorCode(_: domain.Player): string {
    return "";
  }
}
