import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { RatingConverterFactory } from "src/lib/rating/RatingConverter";
import { CssClass, type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class PR
  extends AbstractSingleColumn<CommonKey>
  implements ISummaryColumn
{
  constructor(
    private userConfig: domain.UserConfig,
    private _category: StatsCategory,
  ) {
    super();
  }

  displayKey(): CommonKey {
    return "pr";
  }

  minDisplayName(): string {
    return "PR";
  }

  fullDisplayName(): string {
    return "Personal Rating";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays[this._category].pr;
  }

  tdClass(player: domain.Player): string {
    return this.value(player) === -1 ? CssClass.TD_MULTI : CssClass.TD_NUM;
  }

  displayValue(player: domain.Player): string {
    const value = this.value(player);
    if (value === -1) {
      return "N/A";
    }

    return value.toFixed(this.digit());
  }

  textColorCode(player: domain.Player): string {
    return RatingConverterFactory.fromPR(
      this.value(player),
      this.userConfig,
    ).textColorCode();
  }

  value(player: domain.Player): number {
    return toPlayerStats(player, this.userConfig.stats_pattern)[this._category]
      .pr;
  }

  digit(): number {
    return this.userConfig.custom_digit.pr;
  }

  category(): StatsCategory {
    return this._category;
  }
}
