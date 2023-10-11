import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { RatingConverterFactory } from "src/lib/rating/RatingConverter";
import { CssClass, type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class PR
  extends AbstractColumn<CommonKey>
  implements ISingleColumn, ISummaryColumn
{
  constructor(
    svelteComponent: any,
    private config: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super("pr", "PR", "Personal Rating", 1, svelteComponent);
  }

  shouldShowColumn(): boolean {
    return this.config.displays[this.category].pr;
  }

  getTdClass(player: domain.Player): string {
    return this.getValue(player) === -1 ? CssClass.TD_MULTI : CssClass.TD_NUM;
  }

  getDisplayValue(player: domain.Player): string {
    const value = this.getValue(player);
    if (value === -1) {
      return "N/A";
    }

    return value.toFixed(this.getDigit());
  }

  getTextColorCode(player: domain.Player): string {
    return RatingConverterFactory.fromPR(
      this.getValue(player),
      this.config,
    ).getTextColorCode();
  }

  getValue(player: domain.Player): number {
    return toPlayerStats(player, this.config.stats_pattern)[this.category].pr;
  }

  getDigit(): number {
    return this.config.custom_digit.pr;
  }

  getCategory(): StatsCategory {
    return this.category;
  }
}
