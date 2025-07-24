import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { type Rating, RatingfGenerator } from "src/lib/RatingLevel";
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
    return this.getRating(player)?.textColor ?? "";
  }

  getBackgroundColorCode(player: data.Player): string | undefined {
    return this.getRating(player)?.bgColor;
  }

  getCssClass(): string | undefined {
    if (this.category !== "ship") return "text-right";
    return "text-right";
  }

  private getRating(player: data.Player): Rating | undefined {
    if (this.category !== "ship") return undefined;
    const value = this.playerStats(player).ship.damage;

    return RatingfGenerator.fromDamage(value, player.ship_info.avg_damage);
  }
}
