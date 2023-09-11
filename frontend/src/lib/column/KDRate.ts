import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { CssClass, type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class KDRate
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
    return "kd_rate";
  }

  minDisplayName(): string {
    return "K/D";
  }

  fullDisplayName(): string {
    return "キルデス比";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays[this._category].kd_rate;
  }

  tdClass(player: domain.Player): string {
    return CssClass.TD_NUM;
  }

  displayValue(player: domain.Player): string {
    return this.value(player).toFixed(this.digit());
  }

  textColorCode(player: domain.Player): string {
    return "";
  }

  value(player: domain.Player): number {
    return toPlayerStats(player, this.userConfig.stats_pattern)[this._category]
      .kd_rate;
  }

  digit(): number {
    return this.userConfig.custom_digit.kd_rate;
  }

  category(): StatsCategory {
    return this._category;
  }
}
