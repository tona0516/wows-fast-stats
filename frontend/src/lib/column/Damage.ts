import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { RatingConverterFactory } from "src/lib/rating/RatingConverter";
import { CssClass, type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class Damage
  extends AbstractColumn<CommonKey>
  implements ISingleColumn, ISummaryColumn
{
  constructor(
    svelteComponent: any,
    private config: domain.UserConfig,
    private category: StatsCategory,
  ) {
    super("damage", "Dmg", "平均ダメージ", 1, svelteComponent);
  }

  shouldShowColumn(): boolean {
    return this.config.displays[this.category].damage;
  }

  getTdClass(_: domain.Player): string {
    return CssClass.TD_NUM;
  }

  getDisplayValue(player: domain.Player): string {
    return this.getValue(player).toFixed(this.getDigit());
  }

  getTextColorCode(player: domain.Player): string {
    if (this.category !== "ship") return "";
    const value = toPlayerStats(player, this.config.stats_pattern).ship.damage;

    return RatingConverterFactory.fromDamage(
      value,
      player.ship_info.avg_damage,
      this.config,
    ).getTextColorCode();
  }

  getValue(player: domain.Player): number {
    return toPlayerStats(player, this.config.stats_pattern)[this.category]
      .damage;
  }

  getDigit(): number {
    return this.config.custom_digit.damage;
  }

  getCategory(): StatsCategory {
    return this.category;
  }
}
