import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import { type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export class PlatoonRate
  extends AbstractColumn<CommonKey>
  implements ISingleColumn
{
  constructor(
    private config: model.UserConfig,
    private category: StatsCategory,
  ) {
    super("platoon_rate", "プラ率", "分艦隊比率", 1);
  }

  getSvelteComponent() {
    return SingleTableData;
  }

  shouldShowColumn(): boolean {
    return this.config.displays[this.category].platoon_rate;
  }

  getTdClass(_: model.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: model.Player): string {
    const digit = this.config.custom_digit.platoon_rate;
    const value = toPlayerStats(player, this.config.stats_pattern)[
      this.category
    ].platoon_rate;
    return value.toFixed(digit);
  }

  getTextColorCode(_: model.Player): string {
    return "";
  }
}
