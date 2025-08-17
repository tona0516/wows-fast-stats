import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { StatsCategory } from "src/lib/types";
import type { data } from "wailsjs/go/models";

export class KDRate extends AbstractStatsColumn<string> {
  constructor(category: StatsCategory) {
    super("kd_rate", category);
  }

  displayValue(player: data.Player): string {
    const value = this.playerStats(player)[this.category].kd_rate;
    return value.format(this.digit());
  }

  getTableDataComponent() {
    return SingleTableData;
  }

  getCssClass(): string | undefined {
    return "text-right";
  }
}
