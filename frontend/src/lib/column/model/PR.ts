import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { RatingInfo } from "src/lib/RatingLevel";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
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

    return value.toFixed(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }

  textColorCode(player: data.Player): string {
    return (
      RatingInfo.fromPR(this.value(player), this.config.color.skill.text)
        ?.textColorCode ?? ""
    );
  }

  private value(player: data.Player): number {
    return this.playerStats(player)[this.category].pr;
  }
}
