import { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
import type { ISummaryColumn } from "src/lib/column/intetface/ISummaryColumn";
import { RatingConverterFactory } from "src/lib/rating/RatingConverter";
import { CssClass, type CommonKey, type StatsCategory } from "src/lib/types";
import { toPlayerStats } from "src/lib/util";
import type { domain } from "wailsjs/go/models";

export class Damage
  extends AbstractSingleColumn<CommonKey>
  implements ISummaryColumn
{
  constructor(
    private userConfig: domain.UserConfig,
    private _category: StatsCategory,
  ) {
    super();
  }

  displayKey(): CommonKey {
    return "damage";
  }

  minDisplayName(): string {
    return "Dmg";
  }

  fullDisplayName(): string {
    return "平均ダメージ";
  }

  shouldShowColumn(): boolean {
    return this.userConfig.displays[this._category].damage;
  }

  tdClass(player: domain.Player): string {
    return CssClass.TD_NUM;
  }

  displayValue(player: domain.Player): string {
    return this.value(player).toFixed(this.digit());
  }

  textColorCode(player: domain.Player): string {
    if (this._category !== "ship") return "";
    const value = toPlayerStats(player, this.userConfig.stats_pattern).ship
      .damage;

    return RatingConverterFactory.fromDamage(
      value,
      player.ship_info.avg_damage,
      this.userConfig,
    ).textColorCode();
  }

  value(player: domain.Player): number {
    return toPlayerStats(player, this.userConfig.stats_pattern)[this._category]
      .damage;
  }

  digit(): number {
    return this.userConfig.custom_digit.damage;
  }

  category(): StatsCategory {
    return this._category;
  }
}
