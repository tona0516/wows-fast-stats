import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { Color } from "src/lib/Color";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import { RatingfGenerator } from "src/lib/RatingLevel";
import type { StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class PR extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2, category: StatsCategory) {
    super("pr", config, category);
  }

  displayValue(player: data.Player): string {
    const value = this.value(player);
    if (value === -1) {
      return "N/A";
    }

    return value.format(this.digit());
  }

  getTableDataComponent() {
    return SingleTableData;
  }

  textColorCode(player: data.Player): string {
    const rating = RatingfGenerator.fromPR(this.value(player));
    if (!rating) {
      return "";
    }

    return Color.Rating.getFixed(rating.level)?.text || "";
  }

  getCssClass(): string | undefined {
    return "text-right";
  }

  private value(player: data.Player): number {
    return this.playerStats(player)[this.category].pr;
  }
}
