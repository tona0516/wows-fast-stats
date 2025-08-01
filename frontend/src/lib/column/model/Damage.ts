import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { Color } from "src/lib/Color";
import { RatingfGenerator } from "src/lib/RatingLevel";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class Damage extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2, category: StatsCategory) {
    super("damage", config, category);
  }

  displayValue(player: data.Player): string {
    return this.playerStats(player)[this.category].damage.format(this.digit());
  }

  svelteComponent() {
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

    return Color.Rating.getFixed(rating.level);
  }

  getCssClass(): string | undefined {
    if (this.category !== "ship") return "text-right";
    return "text-right";
  }
}
