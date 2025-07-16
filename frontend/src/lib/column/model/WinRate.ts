import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { RatingInfo } from "src/lib/RatingLevel";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class WinRate extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2, category: StatsCategory) {
    super("win_rate", config, category);
  }

  displayValue(player: data.Player): string {
    return `${this.value(player).format(this.digit())}%`;
  }

  svelteComponent() {
    return SingleTableData;
  }

  textColorCode(player: data.Player): string {
    return (
      RatingInfo.fromWinRate(this.value(player), this.config.color.skill.text)
        ?.textColorCode ?? ""
    );
  }

  private value(player: data.Player): number {
    return this.playerStats(player)[this.category].win_rate;
  }
}
