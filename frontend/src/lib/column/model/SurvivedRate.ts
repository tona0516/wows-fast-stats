import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class SurvivedRate extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("survived_rate", category);
  }

  displayValue(player: data.Player): string {
    const sv = this.playerStats(player)[this.category].survived_rate;

    const all = sv.all.format(this.digit());
    const win = sv.all.format(this.digit());
    const lose = sv.all.format(this.digit());

    return `${all}|${win}|${lose}`;
  }

  getTableDataComponent() {
    return SingleTableData;
  }

  getCssClass(): string | undefined {
    return "text-center";
  }
}
