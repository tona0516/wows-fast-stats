import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AppConst } from "src/lib/AppConst";
import { AppFunc } from "src/lib/AppFunc";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class Damage extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("damage", category);
  }

  displayValue(player: data.Player): string {
    return this.playerStats(player)[this.category].damage.toFixed(this.digit());
  }

  getTableDataComponent() {
    return SingleTableData;
  }

  textColorCode(player: data.Player): string {
    if (this.category !== "ship") "";
    const value = this.playerStats(player).ship.damage;

    const ratingInfo = AppFunc.getRatingfromDamage(
      value,
      player.ship_info.avg_damage,
    );
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

  value(_player: data.Player): number {
    throw new Error("Method not implemented.");
  }

  getCssClass(): string | undefined {
    if (this.category !== "ship") return "text-right";
    return "text-right";
  }
}
