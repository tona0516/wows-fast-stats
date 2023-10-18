import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import { type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class Battles
  extends AbstractColumn<CommonKey>
  implements ISingleColumn
{
  constructor(
    private config: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super("battles", "戦闘数", "戦闘数", 1);
  }

  getSvelteComponent() {
    return SingleTableData;
  }

  shouldShowColumn(): boolean {
    return this.config.displays[this.category].battles;
  }

  getTdClass(_: domain.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: domain.Player): string {
    const digit = this.config.custom_digit.battles;
    const value = toPlayerStats(player, this.config.stats_pattern)[
      this.category
    ].battles;
    return value.toFixed(digit);
  }

  getTextColorCode(_: domain.Player): string {
    return "";
  }
}
