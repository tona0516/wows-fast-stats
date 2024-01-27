import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { CssClass } from "src/lib/CssClass";
import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { RatingColorFactory } from "src/lib/rating/RatingColorFactory";
import { type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { model } from "wailsjs/go/models";

export class Damage
  extends AbstractColumn<CommonKey>
  implements ISingleColumn, ISummaryColumn
{
  constructor(
    private config: model.UserConfig,
    private category: StatsCategory,
  ) {
    super("damage", "Dmg", "平均ダメージ", 1);
  }

  getSvelteComponent() {
    return SingleTableData;
  }

  shouldShowColumn(): boolean {
    return this.config.display[this.category].damage;
  }

  getTdClass(_: model.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: model.Player): string {
    return this.getValue(player).toFixed(this.getDigit());
  }

  getTextColorCode(player: model.Player): string {
    if (this.category !== "ship") return "";
    const value = toPlayerStats(player, this.config.stats_pattern).ship.damage;

    return RatingColorFactory.fromDamage(
      value,
      player.ship_info.avg_damage,
      this.config,
    ).getTextColorCode();
  }

  getValue(player: model.Player): number {
    return toPlayerStats(player, this.config.stats_pattern)[this.category]
      .damage;
  }

  getDigit(): number {
    return this.config.digit.damage;
  }

  getCategory(): StatsCategory {
    return this.category;
  }
}
