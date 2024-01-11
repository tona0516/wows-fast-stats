import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { RatingColorFactory } from "src/lib/rating/RatingColorFactory";
import { type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export class PR
  extends AbstractColumn<CommonKey>
  implements ISingleColumn, ISummaryColumn
{
  constructor(
    private config: model.UserConfig,
    private category: StatsCategory,
  ) {
    super("pr", "PR", "Personal Rating", 1);
  }

  getSvelteComponent() {
    return SingleTableData;
  }

  shouldShowColumn(): boolean {
    return this.config.displays[this.category].pr;
  }

  getTdClass(player: model.Player): string {
    return this.getValue(player) === -1 ? CssClass.TD_MULTI : CssClass.TD_NUM;
  }

  getDisplayValue(player: model.Player): string {
    const value = this.getValue(player);
    if (value === -1) {
      return "N/A";
    }

    return value.toFixed(this.getDigit());
  }

  getTextColorCode(player: model.Player): string {
    return RatingColorFactory.fromPR(
      this.getValue(player),
      this.config,
    ).getTextColorCode();
  }

  getValue(player: model.Player): number {
    return toPlayerStats(player, this.config.stats_pattern)[this.category].pr;
  }

  getDigit(): number {
    return this.config.custom_digit.pr;
  }

  getCategory(): StatsCategory {
    return this.category;
  }
}
