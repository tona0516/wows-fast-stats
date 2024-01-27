import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import { type OverallKey } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export class AvgTier
  extends AbstractColumn<OverallKey>
  implements ISingleColumn
{
  constructor(private config: model.UserConfig) {
    super("avg_tier", "平均T", "平均Tier", 1);
  }

  getSvelteComponent() {
    return SingleTableData;
  }

  shouldShowColumn(): boolean {
    return this.config.display.overall.avg_tier;
  }

  getTdClass(_: model.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: model.Player): string {
    const digit = this.config.digit.avg_tier;
    const value = toPlayerStats(player, this.config.stats_pattern).overall
      .avg_tier;
    return value.toFixed(digit);
  }

  getTextColorCode(_: model.Player): string {
    return "";
  }
}
