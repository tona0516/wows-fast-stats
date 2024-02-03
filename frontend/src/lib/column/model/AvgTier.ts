import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
import type { model } from "wailsjs/go/models";

export class AvgTier extends AbstractStatsColumn<string> {
  constructor(config: model.UserConfig) {
    super("avg_tier", 1, config, "overall");
  }

  displayValue(player: model.Player): string {
    const value = this.playerStats(player).overall.avg_tier;
    return value.toFixed(this.digit());
  }

  svelteComponent() {
    return SingleTableData;
  }
}
