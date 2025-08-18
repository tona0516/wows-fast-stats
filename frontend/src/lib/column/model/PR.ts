import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AppConst } from "src/lib/AppConst";
import { AppFunc } from "src/lib/AppFunc";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class PR extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("pr", category);
  }

  displayValue(player: data.Player): string {
    const value = this.value(player);
    if (value === -1) {
      return "N/A";
    }

    return value.toFixed(this.digit());
  }

  getTableDataComponent() {
    return SingleTableData;
  }

  textColorCode(player: data.Player): string {
    const ratingInfo = AppFunc.getRatingfromPR(this.value(player));
    if (!ratingInfo) {
      return "";
    }

    let fixedColor = "";
    const color = AppConst.RATING_COLORS.get(ratingInfo.rating);
    if (color) {
      fixedColor = AppFunc.getFixedColorPair(color).text;
    }

    return fixedColor;
  }

  getCssClass(): string | undefined {
    return "text-right";
  }

  private value(player: data.Player): number {
    return this.playerStats(player)[this.category].pr;
  }
}
