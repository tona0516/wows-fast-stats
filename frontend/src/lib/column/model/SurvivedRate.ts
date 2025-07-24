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
    const win = stats.win_survived_rate.format(this.digit());
    const lose = stats.lose_survived_rate.format(this.digit());

    return `${win}% | ${lose}%`;
  }

  svelteComponent() {
    return SingleTableData;
  }

  getCssClass(): string | undefined {
    return "text-center";
  }
}
