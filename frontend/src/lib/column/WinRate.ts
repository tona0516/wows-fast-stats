import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { RatingConverterFactory } from "src/lib/rating/RatingConverter";
import { CssClass, type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class WinRate
  extends AbstractColumn<CommonKey>
  implements ISingleColumn, ISummaryColumn
{
  constructor(
    svelteComponent: any,
    private config: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super("win_rate", "勝率", "勝率", 1, svelteComponent);
  }

  shouldShowColumn(): boolean {
    return this.config.displays[this.category].win_rate;
  }

  getTdClass(_: domain.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: domain.Player): string {
    return this.getValue(player).toFixed(this.getDigit()) + "%";
  }

  getTextColorCode(player: domain.Player): string {
    return RatingConverterFactory.fromWinRate(
      this.getValue(player),
      this.config,
    ).getTextColorCode();
  }

  getValue(player: domain.Player): number {
    return toPlayerStats(player, this.config.stats_pattern)[this.category]
      .win_rate;
  }

  getDigit(): number {
    return this.config.custom_digit.win_rate;
  }

  getCategory(): StatsCategory {
    return this.category;
  }
}
