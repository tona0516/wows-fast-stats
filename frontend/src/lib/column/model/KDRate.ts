import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import { type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export class KDRate extends AbstractColumn<CommonKey> implements ISingleColumn {
  constructor(
    private config: model.UserConfig,
    private category: StatsCategory,
  ) {
    super("kd_rate", "K/D", "キル/デス比", 1);
  }

  getSvelteComponent() {
    return SingleTableData;
  }

  shouldShowColumn(): boolean {
    return this.config.displays[this.category].kd_rate;
  }

  getTdClass(_: model.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: model.Player): string {
    const value = toPlayerStats(player, this.config.stats_pattern)[
      this.category
    ].kd_rate;
    const digit = this.config.custom_digit.kd_rate;

    return value.toFixed(digit);
  }

  getTextColorCode(_: model.Player): string {
    return "";
  }
}
