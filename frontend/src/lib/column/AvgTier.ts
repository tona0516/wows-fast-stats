import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import { CssClass, type OverallKey } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class AvgTier
  extends AbstractColumn<OverallKey>
  implements ISingleColumn
{
  constructor(
    svelteComponent: any,
    private config: domain.UserConfig,
  ) {
    super("avg_tier", "平均T", "平均Tier", 1, svelteComponent);
  }

  shouldShowColumn(): boolean {
    return this.config.displays.overall.avg_tier;
  }

  getTdClass(_: domain.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: domain.Player): string {
    const digit = this.config.custom_digit.avg_tier;
    const value = toPlayerStats(player, this.config.stats_pattern).overall
      .avg_tier;
    return value.toFixed(digit);
  }

  getTextColorCode(player: domain.Player): string {
    return "";
  }
}
