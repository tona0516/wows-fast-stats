import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class SurvivedRate extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2, category: StatsCategory) {
    super("survived_rate", config, category);
  }

  displayValue(player: data.Player): string {
    const stats = this.playerStats(player)[this.category];
    return `${stats.win_survived_rate.toFixed(
      this.digit(),
    )}% | ${stats.lose_survived_rate.toFixed(this.digit())}%`;
  }

  svelteComponent() {
    return SingleTableData;
  }
}
