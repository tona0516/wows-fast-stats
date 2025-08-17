import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { Color } from "src/lib/Color";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import { RatingfGenerator } from "src/lib/RatingLevel";
import type { StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class Damage extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("damage", category);
  }

  displayValue(player: data.Player): string {
    return this.playerStats(player)[this.category].damage.format(this.digit());
  }

  getTableDataComponent() {
    return SingleTableData;
  }

  textColorCode(player: data.Player): string {
    if (this.category !== "ship") "";
    const value = this.playerStats(player).ship.damage;

    const rating = RatingfGenerator.fromDamage(
      value,
      player.ship_info.avg_damage,
    );
    if (!rating) {
      return "";
    }

    return Color.Rating.getFixed(rating.level)?.text || "";
  }

  getCssClass(): string | undefined {
    if (this.category !== "ship") return "text-right";
    return "text-right";
  }
}
