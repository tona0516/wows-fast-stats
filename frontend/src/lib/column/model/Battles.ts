import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import { type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export class Battles
  extends AbstractColumn<CommonKey>
  implements ISingleColumn
{
  constructor(
    private config: model.UserConfig,
    private category: StatsCategory,
  ) {
    super("battles", "戦闘数", "戦闘数", 1);
  }

  getSvelteComponent() {
    return SingleTableData;
  }

  shouldShowColumn(): boolean {
    return this.config.display[this.category].battles;
  }

  getTdClass(_: model.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: model.Player): string {
    const digit = this.config.digit.battles;
    const value = toPlayerStats(player, this.config.stats_pattern)[
      this.category
    ].battles;
    return value.toFixed(digit);
  }

  getTextColorCode(_: model.Player): string {
    return "";
  }
}
