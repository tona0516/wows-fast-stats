import SingleTableData from "src/component/stats/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { data } from "wailsjs/go/models";

export class AvgTier extends AbstractStatsColumn<string> {
  constructor(config: data.UserConfigV2) {
    super("avg_tier", config, "overall");
  }

  displayValue(player: data.Player): string {
    const value = this.playerStats(player).overall.avg_tier;
    return value.format(this.digit());
  }

  getTableDataComponent() {
    return SingleTableData;
  }

  getCssClass(): string | undefined {
    return "text-right";
  }
}
