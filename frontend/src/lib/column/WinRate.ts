import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { RatingConverterFactory } from "src/lib/rating/RatingConverter";
import {
  CssClass,
  type CommonStatsKey,
  type StatsCategory,
} from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class WinRate
  extends AbstractSingleColumn<CommonStatsKey>
  implements ISummaryColumn
{
  constructor(
    private userConfig: domain.UserConfig,
    private _category: StatsCategory,
  ) {
    super();
  }

  displayKey(): CommonStatsKey {
    return "win_rate";
  }

  minDisplayName(): string {
    return "勝率";
  }

  fullDisplayName(): string {
    return "勝率";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays[this._category].win_rate;
  }

  tdClass(player: domain.Player): string {
    return CssClass.TD_NUM;
  }

  displayValue(player: domain.Player): string {
    return this.value(player).toFixed(this.digit()) + "%";
  }

  textColorCode(player: domain.Player): string {
    return RatingConverterFactory.fromWinRate(
      this.value(player),
      this.userConfig,
    ).textColorCode();
  }

  value(player: domain.Player): number {
    return toPlayerStats(player, this.userConfig.stats_pattern)[this._category]
      .win_rate;
  }

  digit(): number {
    return this.userConfig.custom_digit.win_rate;
  }

  category(): StatsCategory {
    return this._category;
  }
}
