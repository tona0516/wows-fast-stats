import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import { CssClass, type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class SurvivedRate
  extends AbstractColumn<CommonKey>
  implements ISingleColumn
{
  constructor(
    svelteComponent: any,
    private config: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super(
      "survived_rate",
      "生存率(勝|負)",
      "生存率 (勝利|敗北)",
      1,
      svelteComponent,
    );
  }

  shouldShowColumn(): boolean {
    return this.config.displays[this.category].survived_rate;
  }

  getTdClass(_: domain.Player): string {
    return CssClass.TD_MULTI;
  }

  getDisplayValue(player: domain.Player): string {
    const digit = this.config.custom_digit.survived_rate;
    const stats = toPlayerStats(player, this.config.stats_pattern)[
      this.category
    ];
    return `${stats.win_survived_rate.toFixed(
      digit,
    )}% | ${stats.lose_survived_rate.toFixed(digit)}%`;
  }

  getTextColorCode(_: domain.Player): string {
    return "";
  }
}
